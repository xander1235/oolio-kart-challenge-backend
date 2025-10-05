package controllers_test

import (
	"context"
	"github.com/stretchr/testify/mock"
	"oolio.com/kart/dtos/requests"
	"oolio.com/kart/dtos/responses"
	"oolio.com/kart/exceptions/errors"
	"oolio.com/kart/models"
)

// MockProductService is a mock implementation of ProductService
type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) GetProducts(ctx context.Context, offset, limit *int) []*models.Product {
	args := m.Called(ctx, offset, limit)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).([]*models.Product)
}

func (m *MockProductService) GetProductById(ctx context.Context, id int64) (*models.Product, *errors.ErrorDetails) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*errors.ErrorDetails)
	}
	return args.Get(0).(*models.Product), nil
}

// MockOrderService is a mock implementation of OrderService
type MockOrderService struct {
	mock.Mock
}

func (m *MockOrderService) PlaceOrder(ctx context.Context, request *requests.PlaceOrderRequest) (*responses.OrderResponse, *errors.ErrorDetails) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*errors.ErrorDetails)
	}
	return args.Get(0).(*responses.OrderResponse), nil
}
