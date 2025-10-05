package base

import (
	"context"

	"oolio.com/kart/exceptions/errors"
)

type CouponService interface {
	// ValidateCoupon checks if a coupon code is valid
	ValidateCoupon(ctx context.Context, code string) (bool, *errors.ErrorDetails)
}
