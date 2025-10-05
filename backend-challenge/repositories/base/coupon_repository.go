package base

import (
	"context"

	"oolio.com/kart/exceptions/errors"
)

type CouponRepository interface {
	// GetCouponCounts returns the count of invalid and valid coupons
	GetCouponCounts(ctx context.Context) (invalidCount int64, validCount int64, err *errors.ErrorDetails)

	// GetCouponsByFileCount returns coupon codes filtered by file count condition
	GetCouponsByFileCount(ctx context.Context, fileCountCondition string) ([]string, *errors.ErrorDetails)

	// GetCouponFileCount returns the file count for a specific coupon code
	GetCouponFileCount(ctx context.Context, code string) (fileCount int, found bool, err *errors.ErrorDetails)
}
