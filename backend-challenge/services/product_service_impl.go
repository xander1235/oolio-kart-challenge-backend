package services

import (
	"context"
	"oolio.com/kart/exceptions/errors"
	"oolio.com/kart/models"
	"oolio.com/kart/repositories/base"
)

type ProductServiceImpl struct {
	productRepository base.ProductRepository
}

// NewProductServiceImpl creates a new instance of ProductServiceImpl
func NewProductServiceImpl(productRepository base.ProductRepository) *ProductServiceImpl {
	return &ProductServiceImpl{
		productRepository: productRepository,
	}
}

// GetProducts Retrieves a list of products from the database
func (p *ProductServiceImpl) GetProducts(ctx context.Context, offset, limit *int) []*models.Product {
	return p.productRepository.ListProducts(ctx, offset, limit)
}

// GetProductById Retrieves a product by its ID from the database
func (p *ProductServiceImpl) GetProductById(ctx context.Context, id int64) (*models.Product, *errors.ErrorDetails) {
	return p.productRepository.GetById(ctx, id)
}
