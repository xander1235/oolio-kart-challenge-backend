package base

import (
	"context"

	"oolio.com/kart/exceptions/errors"
	"oolio.com/kart/models"
)

type OrderRepository interface {
	// CreateOrder creates a new order in the database
	CreateOrder(ctx context.Context, order *models.Order, items []models.OrderItem) *errors.ErrorDetails
}
