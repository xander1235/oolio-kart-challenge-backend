package services_test

import (
	"context"
	"github.com/stretchr/testify/mock"
	"oolio.com/kart/exceptions/errors"
	"oolio.com/kart/models"
)

// MockProductRepository is a mock implementation of ProductRepository
type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) ListProducts(ctx context.Context, offset, limit *int) []*models.Product {
	args := m.Called(ctx, offset, limit)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).([]*models.Product)
}

func (m *MockProductRepository) GetById(ctx context.Context, id int64) (*models.Product, *errors.ErrorDetails) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*errors.ErrorDetails)
	}
	return args.Get(0).(*models.Product), nil
}

func (m *MockProductRepository) Save(ctx context.Context, product *models.Product) *errors.ErrorDetails {
	args := m.Called(ctx, product)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*errors.ErrorDetails)
}

func (m *MockProductRepository) Update(ctx context.Context, product *models.Product) *errors.ErrorDetails {
	args := m.Called(ctx, product)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*errors.ErrorDetails)
}

func (m *MockProductRepository) GetByIds(ctx context.Context, ids []int64) ([]*models.Product, *errors.ErrorDetails) {
	args := m.Called(ctx, ids)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*errors.ErrorDetails)
	}
	return args.Get(0).([]*models.Product), nil
}

// MockOrderRepository is a mock implementation of OrderRepository
type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) CreateOrder(ctx context.Context, order *models.Order, items []models.OrderItem) *errors.ErrorDetails {
	args := m.Called(ctx, order, items)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*errors.ErrorDetails)
}

// MockCouponRepository is a mock implementation of CouponRepository
type MockCouponRepository struct {
	mock.Mock
}

func (m *MockCouponRepository) GetCouponCounts(ctx context.Context) (invalidCount int64, validCount int64, err *errors.ErrorDetails) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Get(1).(int64), nil
}

func (m *MockCouponRepository) GetCouponsByFileCount(ctx context.Context, fileCountCondition string) ([]string, *errors.ErrorDetails) {
	args := m.Called(ctx, fileCountCondition)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*errors.ErrorDetails)
	}
	return args.Get(0).([]string), nil
}

func (m *MockCouponRepository) GetCouponFileCount(ctx context.Context, code string) (fileCount int, found bool, err *errors.ErrorDetails) {
	args := m.Called(ctx, code)
	return args.Int(0), args.Bool(1), nil
}
