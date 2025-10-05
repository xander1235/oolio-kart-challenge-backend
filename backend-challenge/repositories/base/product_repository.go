package base

import (
	"context"
	"oolio.com/kart/exceptions/errors"
	"oolio.com/kart/models"
)

type ProductRepository interface {
	// Save saves a new product to the database
	Save(ctx context.Context, product *models.Product) *errors.ErrorDetails

	// Update updates an existing product in the database
	Update(ctx context.Context, product *models.Product) *errors.ErrorDetails

	// GetById retrieves a product by its ID from the database
	GetById(ctx context.Context, id int64) (*models.Product, *errors.ErrorDetails)

	// ListProducts retrieves a list of products from the database
	ListProducts(ctx context.Context, limit, offset *int) []*models.Product
}
