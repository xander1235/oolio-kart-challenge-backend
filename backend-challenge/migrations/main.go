package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSchema   string
	DBSSLMode  string

	CouponSource         string
	CouponDataDir        string
	CouponS3BaseURL      string
	CouponForceMigration bool

	ProductDataFile string
}

func main() {
	migrationType := flag.String("type", "all", "Migration type: coupon, product, or all")
	envFile := flag.String("env", ".env", "Path to .env file")
	flag.Parse()

	if err := godotenv.Load(*envFile); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	config := loadConfig()

	ctx := context.Background()
	pool, err := createDBPool(ctx, config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	log.Println("Connected to database successfully")

	switch *migrationType {
	case "coupon":
		if err := runCouponMigration(ctx, pool, config); err != nil {
			log.Fatalf("Coupon migration failed: %v", err)
		}
	case "product":
		if err := runProductMigration(ctx, pool, config); err != nil {
			log.Fatalf("Product migration failed: %v", err)
		}
	case "all":
		log.Println("Running all migrations...")
		if err := runProductMigration(ctx, pool, config); err != nil {
			log.Fatalf("Product migration failed: %v", err)
		}
		if err := runCouponMigration(ctx, pool, config); err != nil {
			log.Fatalf("Coupon migration failed: %v", err)
		}
	default:
		log.Fatalf("Invalid migration type: %s. Use: coupon, product, or all", *migrationType)
	}

	log.Println("All migrations completed successfully!")
}

func loadConfig() *Config {
	forceMigration, _ := strconv.ParseBool(getEnv("COUPON_FORCE_MIGRATION", "false"))

	return &Config{
		DBHost:               getEnv("DB_HOST", "localhost"),
		DBPort:               getEnv("DB_PORT", "5432"),
		DBUser:               getEnv("DB_USER", ""),
		DBPassword:           getEnv("DB_PASSWORD", ""),
		DBName:               getEnv("DB_NAME", "postgres"),
		DBSchema:             getEnv("DB_SCHEMA", "kart"),
		DBSSLMode:            getEnv("DB_SSL_MODE", "disable"),
		CouponSource:         getEnv("COUPON_SOURCE", "local"),
		CouponDataDir:        getEnv("COUPON_DATA_DIR", "../data"),
		CouponS3BaseURL:      getEnv("COUPON_S3_BASE_URL", "https://orderfoodonline-files.s3.ap-southeast-2.amazonaws.com"),
		CouponForceMigration: forceMigration,
		ProductDataFile:      getEnv("PRODUCT_DATA_FILE", "../data/product.json"),
	}
}

func createDBPool(ctx context.Context, config *Config) (*pgxpool.Pool, error) {
	var connStr string
	if config.DBUser != "" && config.DBPassword != "" {
		connStr = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=%s&search_path=%s",
			config.DBUser,
			config.DBPassword,
			config.DBHost,
			config.DBPort,
			config.DBName,
			config.DBSSLMode,
			config.DBSchema,
		)
	} else if config.DBUser != "" && config.DBPassword == "" {
		connStr = fmt.Sprintf(
			"postgres://%s@%s:%s/%s?sslmode=%s&search_path=%s",
			config.DBUser,
			config.DBHost,
			config.DBPort,
			config.DBName,
			config.DBSSLMode,
			config.DBSchema,
		)
	} else {
		connStr = fmt.Sprintf(
			"postgres://%s:%s/%s?sslmode=%s&search_path=%s",
			config.DBHost,
			config.DBPort,
			config.DBName,
			config.DBSSLMode,
			config.DBSchema,
		)
	}

	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	poolConfig.MaxConns = 10
	poolConfig.MinConns = 2
	poolConfig.MaxConnLifetime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}

func runCouponMigration(ctx context.Context, pool *pgxpool.Pool, config *Config) error {
	log.Println("Starting coupon migration...")

	couponConfig := &CouponConfiguration{
		Source:         config.CouponSource,
		DataDir:        config.CouponDataDir,
		S3BaseURL:      config.CouponS3BaseURL,
		ForceMigration: config.CouponForceMigration,
	}

	migration := NewCouponMigration(pool, couponConfig)
	if err := migration.Run(ctx); err != nil {
		return err
	}

	// Analyze distribution
	invalidCount, validCount, err := migration.AnalyzeDistribution(ctx)
	if err != nil {
		log.Printf("Warning: Could not analyze distribution: %v", err)
	} else {
		log.Printf("Coupon distribution - Invalid: %d, Valid: %d", invalidCount, validCount)
	}

	return nil
}

func runProductMigration(ctx context.Context, pool *pgxpool.Pool, config *Config) error {
	log.Println("Starting product migration...")

	migration := NewProductMigration(pool, config.ProductDataFile)
	return migration.Run(ctx)
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
