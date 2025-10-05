package repositories

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"oolio.com/kart/configs"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"oolio.com/kart/exceptions"
	"oolio.com/kart/exceptions/errors"
	"oolio.com/kart/repositories/base"
)

type CouponRepositoryImpl struct {
	pool *pgxpool.Pool
}

// NewCouponRepositoryImpl creates a new instance of CouponRepositoryImpl
func NewCouponRepositoryImpl(pool *pgxpool.Pool) base.CouponRepository {
	return &CouponRepositoryImpl{pool: pool}
}

// GetCouponCounts returns the count of invalid and valid coupons
func (r *CouponRepositoryImpl) GetCouponCounts(ctx context.Context) (invalidCount int64, validCount int64, err *errors.ErrorDetails) {
	err2 := r.pool.QueryRow(ctx, `
		SELECT 
			COUNT(*) FILTER (WHERE file_count = 1) AS invalid_count,
			COUNT(*) FILTER (WHERE file_count >= 2) AS valid_count
		FROM coupons
	`).Scan(&invalidCount, &validCount)

	if err2 != nil {
		configs.Logger.Error("failed to count coupons", zap.Error(err2))
		return 0, 0, exceptions.GenericException("failed to count coupons", http.StatusInternalServerError)
	}

	return invalidCount, validCount, nil
}

// GetCouponsByFileCount returns coupon codes filtered by file count condition
// fileCountCondition should be either "= 1" for invalid or ">= 2" for valid
func (r *CouponRepositoryImpl) GetCouponsByFileCount(ctx context.Context, fileCountCondition string) ([]string, *errors.ErrorDetails) {
	query := fmt.Sprintf("SELECT code FROM coupons WHERE file_count %s", fileCountCondition)

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		configs.Logger.Error("failed to query coupons", zap.Error(err))
		return nil, exceptions.GenericException("failed to query coupons", http.StatusInternalServerError)
	}
	defer rows.Close()

	var codes []string
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			continue
		}
		codes = append(codes, code)
	}

	if err := rows.Err(); err != nil {
		configs.Logger.Error("error reading coupons", zap.Error(err))
		return nil, exceptions.GenericException("error reading coupons", http.StatusInternalServerError)
	}

	return codes, nil
}

// GetCouponFileCount returns the file count for a specific coupon code
func (r *CouponRepositoryImpl) GetCouponFileCount(ctx context.Context, code string) (fileCount int, found bool, err *errors.ErrorDetails) {
	err2 := r.pool.QueryRow(ctx,
		"SELECT file_count FROM coupons WHERE code = $1",
		code,
	).Scan(&fileCount)

	if err2 == pgx.ErrNoRows {
		configs.Logger.Error("coupon not found", zap.Error(err2))
		return 0, false, nil
	}
	if err2 != nil {
		configs.Logger.Error("failed to get coupon file count", zap.Error(err2))
		return 0, false, exceptions.GenericException("failed to get coupon file count", http.StatusInternalServerError)
	}

	return fileCount, true, nil
}
