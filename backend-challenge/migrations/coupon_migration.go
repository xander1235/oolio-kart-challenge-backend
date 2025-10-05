package main

import (
	"bufio"
	"compress/gzip"
	"context"
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	NumPartitions = 16
	BatchSize     = 10000
	MaxWorkers    = 8
)

type CouponConfiguration struct {
	Source         string
	DataDir        string
	S3BaseURL      string
	ForceMigration bool
}

type CouponMigration struct {
	pool   *pgxpool.Pool
	config *CouponConfiguration
}

func NewCouponMigration(pool *pgxpool.Pool, config *CouponConfiguration) *CouponMigration {
	return &CouponMigration{
		pool:   pool,
		config: config,
	}
}

func (cl *CouponMigration) Run(ctx context.Context) error {
	needsMigration, err := cl.checkMigrationNeeded(ctx)
	if err != nil {
		return fmt.Errorf("failed to check migration status: %w", err)
	}

	if !needsMigration {
		log.Println("Coupons already loaded, skipping migration")
		return nil
	}

	startTime := time.Now()

	if err := cl.createStagingTable(ctx); err != nil {
		return fmt.Errorf("failed to create staging table: %w", err)
	}

	_, _ = cl.pool.Exec(ctx, "ALTER TABLE coupon_staging SET (autovacuum_enabled = false)")
	_, _ = cl.pool.Exec(ctx, "ALTER TABLE coupons SET (autovacuum_enabled = false)")

	files := []string{"couponbase1.gz", "couponbase2.gz", "couponbase3.gz"}
	errChan := make(chan error, len(files))

	for _, filename := range files {
		go func(fname string) {
			fileSource := strings.TrimSuffix(fname, ".gz")
			log.Printf("Loading coupons from file: %s", fname)

			if err := cl.loadFromFile(ctx, fname, fileSource); err != nil {
				errChan <- fmt.Errorf("failed to load %s: %w", fname, err)
				return
			}
			errChan <- nil
		}(filename)
	}

	for i := 0; i < len(files); i++ {
		if err := <-errChan; err != nil {
			return err
		}
	}

	if err := cl.aggregateCouponsParallel(ctx); err != nil {
		return fmt.Errorf("failed to aggregate coupons: %w", err)
	}

	if err := cl.dropStagingTable(ctx); err != nil {
		log.Printf("Warning: Failed to drop staging table: %v", err)
	}

	_, _ = cl.pool.Exec(ctx, "ALTER TABLE coupons SET (autovacuum_enabled = true)")
	log.Printf("Coupon migration completed in %v", time.Since(startTime))

	return nil
}

func (cl *CouponMigration) checkMigrationNeeded(ctx context.Context) (bool, error) {
	if cl.config.ForceMigration {
		log.Println("Force migration enabled, will reload coupons")
		_, err := cl.pool.Exec(ctx, "TRUNCATE TABLE coupons")
		return true, err
	}

	var count int64
	err := cl.pool.QueryRow(ctx, "SELECT COUNT(*) FROM coupons").Scan(&count)
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

func (cl *CouponMigration) createStagingTable(ctx context.Context) error {
	_, _ = cl.pool.Exec(ctx, "DROP TABLE IF EXISTS coupon_staging")

	query := `
		CREATE UNLOGGED TABLE coupon_staging (
			code         VARCHAR(10) NOT NULL,
			file_source  VARCHAR(20) NOT NULL,
			partition_id INT NOT NULL
		)
	`
	if _, err := cl.pool.Exec(ctx, query); err != nil {
		return err
	}

	_, err := cl.pool.Exec(ctx, "CREATE INDEX idx_staging_partition ON coupon_staging(partition_id)")
	return err
}

func (cl *CouponMigration) loadFromFile(ctx context.Context, filename, fileSource string) error {
	var reader io.ReadCloser
	var err error

	if cl.config.Source == "s3" {
		reader, err = cl.openS3Stream(filename)
	} else {
		reader, err = cl.openLocalFile(filename)
	}

	if err != nil {
		return err
	}
	defer reader.Close()

	gzReader, err := gzip.NewReader(reader)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzReader.Close()

	scanner := bufio.NewScanner(gzReader)
	batch := make([][]any, 0, 10000)
	totalProcessed := 0

	for scanner.Scan() {
		code := strings.TrimSpace(scanner.Text())

		if len(code) < 8 || len(code) > 10 {
			continue
		}

		partitionID := cl.hashPartition(code)
		batch = append(batch, []any{code, fileSource, partitionID})

		if len(batch) >= 10000 {
			if err := cl.batchInsert(ctx, batch); err != nil {
				return err
			}
			totalProcessed += len(batch)
			batch = batch[:0]

			if totalProcessed%1000000 == 0 {
				log.Printf("Progress [%s]: %d rows processed", filename, totalProcessed)
			}
		}
	}

	if len(batch) > 0 {
		if err := cl.batchInsert(ctx, batch); err != nil {
			return err
		}
		totalProcessed += len(batch)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %w", err)
	}

	log.Printf("Completed loading file %s: %d rows", filename, totalProcessed)
	return nil
}

func (cl *CouponMigration) openS3Stream(filename string) (io.ReadCloser, error) {
	url := fmt.Sprintf("%s/%s", cl.config.S3BaseURL, filename)

	client := &http.Client{
		Timeout: 30 * time.Minute,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from S3: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("bad status from S3: %s", resp.Status)
	}

	return resp.Body, nil
}

func (cl *CouponMigration) openLocalFile(filename string) (io.ReadCloser, error) {
	filePath := filepath.Join(cl.config.DataDir, filename)
	return os.Open(filePath)
}

func (cl *CouponMigration) batchInsert(ctx context.Context, batch [][]any) error {
	_, err := cl.pool.CopyFrom(
		ctx,
		pgx.Identifier{"coupon_staging"},
		[]string{"code", "file_source", "partition_id"},
		pgx.CopyFromRows(batch),
	)
	return err
}

func (cl *CouponMigration) hashPartition(code string) int {
	h := fnv.New32a()
	h.Write([]byte(code))
	return int(h.Sum32() % NumPartitions)
}

func (cl *CouponMigration) aggregateCouponsParallel(ctx context.Context) error {
	log.Printf("Aggregating coupons in parallel (partitions: %d, workers: %d)", NumPartitions, MaxWorkers)

	var wg sync.WaitGroup
	errChan := make(chan error, NumPartitions)
	semaphore := make(chan struct{}, MaxWorkers)

	for partitionID := 0; partitionID < NumPartitions; partitionID++ {
		wg.Add(1)
		go func(pid int) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			if err := cl.processPartition(ctx, pid); err != nil {
				errChan <- fmt.Errorf("partition %d failed: %w", pid, err)
				return
			}
		}(partitionID)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	log.Println("All partitions processed successfully")
	return nil
}

func (cl *CouponMigration) processPartition(ctx context.Context, partitionID int) error {
	conn, err := cl.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	optimizations := `
		SET LOCAL work_mem = '1GB';
		SET LOCAL maintenance_work_mem = '2GB';
		SET LOCAL temp_buffers = '1GB';
		SET LOCAL synchronous_commit = OFF;
		SET LOCAL enable_nestloop = OFF;
		SET LOCAL random_page_cost = 1.1;
	`
	if _, err := conn.Exec(ctx, optimizations); err != nil {
		log.Printf("Warning: Failed to apply optimizations for partition %d: %v", partitionID, err)
	}

	query := `
		INSERT INTO coupons (code, file_sources, file_count)
		SELECT 
			code,
			ARRAY_AGG(DISTINCT file_source ORDER BY file_source) as file_sources,
			COUNT(DISTINCT file_source) as file_count
		FROM coupon_staging
		WHERE partition_id = $1
		GROUP BY code
		ON CONFLICT (code) DO NOTHING
	`

	result, err := conn.Exec(ctx, query, partitionID)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	log.Printf("Partition %d processed: %d unique coupons", partitionID, rowsAffected)

	return nil
}

func (cl *CouponMigration) dropStagingTable(ctx context.Context) error {
	_, err := cl.pool.Exec(ctx, "DROP TABLE IF EXISTS coupon_staging")
	return err
}

func (cl *CouponMigration) AnalyzeDistribution(ctx context.Context) (invalidCount, validCount int64, err error) {
	err = cl.pool.QueryRow(ctx,
		"SELECT COUNT(*) FROM coupons WHERE file_count = 1",
	).Scan(&invalidCount)
	if err != nil {
		return 0, 0, err
	}

	err = cl.pool.QueryRow(ctx,
		"SELECT COUNT(*) FROM coupons WHERE file_count >= 2",
	).Scan(&validCount)
	if err != nil {
		return 0, 0, err
	}

	return invalidCount, validCount, nil
}
