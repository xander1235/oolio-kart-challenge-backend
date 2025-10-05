package base

import (
	"context"
	"oolio.com/kart/exceptions/errors"
	"oolio.com/kart/models"
)

type ProductService interface {
	// GetProducts retrieves a list of products from the database
	GetProducts(ctx context.Context, offset, limit *int) []*models.Product

	// GetProductById retrieves a product by its ID from the database
	GetProductById(ctx context.Context, id int64) (*models.Product, *errors.ErrorDetails)
}
