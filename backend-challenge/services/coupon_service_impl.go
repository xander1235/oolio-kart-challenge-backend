package services

import (
	"context"
	"go.uber.org/zap"
	"oolio.com/kart/repositories/base"

	"oolio.com/kart/configs"
	"oolio.com/kart/exceptions/errors"
)

var CouponServiceImpl *couponServiceImpl

type couponServiceImpl struct {
	validator  *CouponValidator
	couponRepo base.CouponRepository
}

// InitializeCouponService initializes the coupon service (loads Bloom filter from database)
// Note: Migrations should be run separately before calling this
func InitializeCouponService(couponRepo base.CouponRepository) *errors.ErrorDetails {
	configs.Logger.Info("Initializing coupon service...")
	validator, err := NewCouponValidator(couponRepo)
	if err != nil {
		configs.Logger.Error("Failed to initialize coupon service", zap.Any("error", err))
		return err
	}

	CouponServiceImpl = &couponServiceImpl{
		validator:  validator,
		couponRepo: couponRepo,
	}

	configs.Logger.Info("Coupon service initialized successfully")
	return nil
}

// ValidateCoupon validates if a coupon code is valid
func (s *couponServiceImpl) ValidateCoupon(ctx context.Context, code string) (bool, *errors.ErrorDetails) {
	return s.validator.ValidateCoupon(ctx, code)
}
