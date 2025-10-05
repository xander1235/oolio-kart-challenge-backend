package base

import (
	"context"

	"oolio.com/kart/dtos/requests"
	"oolio.com/kart/dtos/responses"
	"oolio.com/kart/exceptions/errors"
)

type OrderService interface {
	// PlaceOrder places a new order
	PlaceOrder(ctx context.Context, request *requests.PlaceOrderRequest) (*responses.OrderResponse, *errors.ErrorDetails)
}
