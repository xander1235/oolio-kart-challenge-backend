package services

import (
	"context"
	"github.com/bits-and-blooms/bloom/v3"
	"go.uber.org/zap"
	"oolio.com/kart/configs"
	"oolio.com/kart/exceptions/errors"
	"oolio.com/kart/repositories/base"
)

type CouponValidator struct {
	bloomFilter    *bloom.BloomFilter
	filterStrategy string
	repository     base.CouponRepository
}

// NewCouponValidator creates a new coupon validator and loads the bloom filter
func NewCouponValidator(repository base.CouponRepository) (*CouponValidator, *errors.ErrorDetails) {
	ctx := context.Background()

	invalidCount, validCount, err := repository.GetCouponCounts(ctx)
	if err != nil {
		configs.Logger.Error("failed to get coupon counts", zap.Any("error", err))
		return nil, err
	}

	var strategy string
	var fileCountCondition string
	var count int64

	if invalidCount < validCount {
		strategy = "negative"
		fileCountCondition = "= 1"
		count = invalidCount
		configs.Logger.Info("Using NEGATIVE Bloom filter (storing invalid coupons)",
			zap.Int64("count", count),
		)
	} else {
		strategy = "positive"
		fileCountCondition = ">= 2"
		count = validCount
		configs.Logger.Info("Using POSITIVE Bloom filter (storing valid coupons)",
			zap.Int64("count", count),
		)
	}

	if count == 0 {
		configs.Logger.Warn("no coupons found in database, all coupons will be invalid")
		// Create an empty bloom filter - all coupons will be rejected
		bf := bloom.NewWithEstimates(1, 0.001) // Minimum size
		return &CouponValidator{
			bloomFilter:    bf,
			filterStrategy: strategy,
			repository:     repository,
		}, nil
	}

	bf := bloom.NewWithEstimates(uint(count), 0.001)

	codes, err := repository.GetCouponsByFileCount(ctx, fileCountCondition)
	if err != nil {
		configs.Logger.Error("failed to get coupons by file count", zap.Any("error", err))
		return nil, err
	}

	loaded := 0
	for _, code := range codes {
		bf.Add([]byte(code))
		loaded++

		if loaded%1000000 == 0 {
			configs.Logger.Info("Loading Bloom filter", zap.Int("loaded", loaded))
		}
	}

	configs.Logger.Info("Bloom filter loaded successfully",
		zap.String("strategy", strategy),
		zap.Int("loaded", loaded),
		zap.Int64("total_coupons", invalidCount+validCount),
		zap.Int64("valid_coupons", validCount),
		zap.Int64("invalid_coupons", invalidCount),
		zap.Float64("valid_percentage", float64(validCount)/float64(invalidCount+validCount)*100),
	)

	return &CouponValidator{
		bloomFilter:    bf,
		filterStrategy: strategy,
		repository:     repository,
	}, nil
}

// ValidateCoupon checks if a coupon code is valid
func (v *CouponValidator) ValidateCoupon(ctx context.Context, code string) (bool, *errors.ErrorDetails) {
	// Format validation
	if len(code) < 8 || len(code) > 10 {
		return false, nil
	}

	inBloom := v.bloomFilter.Test([]byte(code))

	if v.filterStrategy == "negative" {
		// Bloom filter stores INVALID coupons (file_count = 1)
		if !inBloom {
			// Not in bloom → definitely VALID (no DB hit needed!)
			return true, nil
		}
		// In bloom → invalid OR false positive → check DB
		return v.checkDBIsValid(ctx, code)

	} else {
		// Bloom filter stores VALID coupons (file_count >= 2)
		if !inBloom {
			// Not in bloom → definitely INVALID
			return false, nil
		}
		// In bloom → valid OR false positive → check DB
		return v.checkDBIsValid(ctx, code)
	}
}

func (v *CouponValidator) checkDBIsValid(ctx context.Context, code string) (bool, *errors.ErrorDetails) {
	fileCount, found, err := v.repository.GetCouponFileCount(ctx, code)
	if err != nil {
		return false, err
	}
	if !found {
		return false, nil
	}

	return fileCount >= 2, nil
}
