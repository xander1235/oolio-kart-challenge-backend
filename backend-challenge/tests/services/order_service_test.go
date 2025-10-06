package services_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"oolio.com/kart/dtos/requests"
	"oolio.com/kart/exceptions/errors"
	"oolio.com/kart/models"
	"oolio.com/kart/services"
	"testing"
)

// TestOrderService_PlaceOrder_Success tests the PlaceOrder method of the OrderService
func TestOrderService_PlaceOrder_Success(t *testing.T) {
	mockCouponRepo := new(MockCouponRepository)
	mockCouponRepo.On("GetCouponCounts", mock.Anything).Return(int64(100), int64(2), nil)
	mockCouponRepo.On("GetCouponsByFileCount", mock.Anything, ">= 2").Return([]string{"SAVE1000", "DISCOUNT50"}, nil)

	err := services.InitializeCouponService(mockCouponRepo)
	assert.Nil(t, err)

	mockOrderRepo := new(MockOrderRepository)
	mockProductRepo := new(MockProductRepository)
	service := services.NewOrderServiceImpl(mockOrderRepo, mockProductRepo, services.CouponServiceImpl)

	quantity := 2
	request := &requests.PlaceOrderRequest{
		Items: []requests.OrderItemRequest{
			{ProductId: "1", Quantity: &quantity},
		},
	}

	mockProducts := []*models.Product{
		{
			Id:       1,
			Name:     "Margherita Pizza",
			Price:    12.99,
			Category: "Pizza",
			Status:   "available",
		},
	}

	mockProductRepo.On("GetByIds", mock.Anything, []int64{1}).Return(mockProducts, nil)
	mockOrderRepo.On("CreateOrder", mock.Anything, mock.AnythingOfType("*models.Order"), mock.AnythingOfType("[]models.OrderItem")).
		Run(func(args mock.Arguments) {
			order := args.Get(1).(*models.Order)
			order.Id = "550e8400-e29b-41d4-a716-446655440000"
		}).
		Return(nil)

	result, err := service.PlaceOrder(context.Background(), request)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.Id)
	assert.Len(t, result.Items, 1)
	assert.Equal(t, 2, result.Items[0].Quantity)

	mockProductRepo.AssertExpectations(t)
	mockOrderRepo.AssertExpectations(t)
}

// TestOrderService_PlaceOrder_WithValidCoupon tests the PlaceOrder method with actual coupon validation logic
func TestOrderService_PlaceOrder_WithValidCoupon(t *testing.T) {
	mockOrderRepo := new(MockOrderRepository)
	mockProductRepo := new(MockProductRepository)
	mockCouponRepo := new(MockCouponRepository)

	mockCouponRepo.On("GetCouponCounts", mock.Anything).Return(int64(100), int64(2), nil)
	mockCouponRepo.On("GetCouponsByFileCount", mock.Anything, ">= 2").Return([]string{"SAVE1000", "DISCOUNT50"}, nil)
	mockCouponRepo.On("GetCouponFileCount", mock.Anything, "SAVE1000").Return(2, true, nil)

	err := services.InitializeCouponService(mockCouponRepo)
	assert.Nil(t, err)

	service := services.NewOrderServiceImpl(mockOrderRepo, mockProductRepo, services.CouponServiceImpl)

	quantity := 2
	request := &requests.PlaceOrderRequest{
		CouponCode: "SAVE1000",
		Items: []requests.OrderItemRequest{
			{ProductId: "1", Quantity: &quantity},
		},
	}

	mockProducts := []*models.Product{
		{
			Id:       1,
			Name:     "Margherita Pizza",
			Price:    12.99,
			Category: "Pizza",
			Status:   "available",
		},
	}

	mockProductRepo.On("GetByIds", mock.Anything, []int64{1}).Return(mockProducts, nil)
	mockOrderRepo.On("CreateOrder", mock.Anything, mock.AnythingOfType("*models.Order"), mock.AnythingOfType("[]models.OrderItem")).
		Run(func(args mock.Arguments) {
			order := args.Get(1).(*models.Order)
			order.Id = "550e8400-e29b-41d4-a716-446655440001"
		}).
		Return(nil)

	result, errDetails := service.PlaceOrder(context.Background(), request)

	assert.Nil(t, errDetails)
	assert.NotNil(t, result)
	assert.Equal(t, "SAVE1000", result.CouponCode)
	assert.NotEmpty(t, result.Id)

	mockCouponRepo.AssertExpectations(t)
	mockProductRepo.AssertExpectations(t)
	mockOrderRepo.AssertExpectations(t)
}

// TestOrderService_PlaceOrder_InvalidCoupon tests with actual coupon validation logic
func TestOrderService_PlaceOrder_InvalidCoupon(t *testing.T) {
	mockOrderRepo := new(MockOrderRepository)
	mockProductRepo := new(MockProductRepository)
	mockCouponRepo := new(MockCouponRepository)

	mockCouponRepo.On("GetCouponCounts", mock.Anything).Return(int64(100), int64(2), nil)
	mockCouponRepo.On("GetCouponsByFileCount", mock.Anything, ">= 2").Return([]string{"SAVE1000", "DISCOUNT50"}, nil)

	err := services.InitializeCouponService(mockCouponRepo)
	assert.Nil(t, err)

	service := services.NewOrderServiceImpl(mockOrderRepo, mockProductRepo, services.CouponServiceImpl)

	quantity := 2
	request := &requests.PlaceOrderRequest{
		CouponCode: "INVALID123",
		Items: []requests.OrderItemRequest{
			{ProductId: "1", Quantity: &quantity},
		},
	}

	result, errDetails := service.PlaceOrder(context.Background(), request)

	assert.Nil(t, result)
	assert.NotNil(t, errDetails)
	assert.Equal(t, http.StatusUnprocessableEntity, errDetails.ErrorCode)
	assert.Contains(t, errDetails.Message, "invalid coupon code")

	mockCouponRepo.AssertExpectations(t)
}

// TestOrderService_PlaceOrder_InvalidCouponFormat tests coupon with invalid format
func TestOrderService_PlaceOrder_InvalidCouponFormat(t *testing.T) {
	mockOrderRepo := new(MockOrderRepository)
	mockProductRepo := new(MockProductRepository)
	mockCouponRepo := new(MockCouponRepository)

	mockCouponRepo.On("GetCouponCounts", mock.Anything).Return(int64(100), int64(1), nil)
	mockCouponRepo.On("GetCouponsByFileCount", mock.Anything, ">= 2").Return([]string{"SAVE1000"}, nil)

	err := services.InitializeCouponService(mockCouponRepo)
	assert.Nil(t, err)

	service := services.NewOrderServiceImpl(mockOrderRepo, mockProductRepo, services.CouponServiceImpl)

	quantity := 2
	request := &requests.PlaceOrderRequest{
		CouponCode: "SHORT",
		Items: []requests.OrderItemRequest{
			{ProductId: "1", Quantity: &quantity},
		},
	}

	result, errDetails := service.PlaceOrder(context.Background(), request)

	assert.Nil(t, result)
	assert.NotNil(t, errDetails)
	assert.Equal(t, http.StatusUnprocessableEntity, errDetails.ErrorCode)

	mockCouponRepo.AssertExpectations(t)
}

// TestOrderService_PlaceOrder_ProductNotFound tests the PlaceOrder method with a product that does not exist
func TestOrderService_PlaceOrder_ProductNotFound(t *testing.T) {
	mockCouponRepo := new(MockCouponRepository)
	mockCouponRepo.On("GetCouponCounts", mock.Anything).Return(int64(100), int64(2), nil)
	mockCouponRepo.On("GetCouponsByFileCount", mock.Anything, ">= 2").Return([]string{"SAVE1000", "DISCOUNT50"}, nil)

	err := services.InitializeCouponService(mockCouponRepo)
	assert.Nil(t, err)

	mockOrderRepo := new(MockOrderRepository)
	mockProductRepo := new(MockProductRepository)
	service := services.NewOrderServiceImpl(mockOrderRepo, mockProductRepo, services.CouponServiceImpl)

	quantity := 2
	request := &requests.PlaceOrderRequest{
		Items: []requests.OrderItemRequest{
			{ProductId: "999", Quantity: &quantity},
		},
	}

	mockError := &errors.ErrorDetails{
		ErrorCode: http.StatusInternalServerError,
		Message:   "failed to fetch products",
	}

	mockProductRepo.On("GetByIds", mock.Anything, []int64{999}).Return(nil, mockError)

	result, err := service.PlaceOrder(context.Background(), request)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusBadRequest, err.ErrorCode)

	mockProductRepo.AssertExpectations(t)
}

// TestOrderService_PlaceOrder_InvalidProductId tests the PlaceOrder method with an invalid product ID
func TestOrderService_PlaceOrder_InvalidProductId(t *testing.T) {
	mockCouponRepo := new(MockCouponRepository)
	mockCouponRepo.On("GetCouponCounts", mock.Anything).Return(int64(100), int64(2), nil)
	mockCouponRepo.On("GetCouponsByFileCount", mock.Anything, ">= 2").Return([]string{"SAVE1000", "DISCOUNT50"}, nil)

	err := services.InitializeCouponService(mockCouponRepo)
	assert.Nil(t, err)

	mockOrderRepo := new(MockOrderRepository)
	mockProductRepo := new(MockProductRepository)
	service := services.NewOrderServiceImpl(mockOrderRepo, mockProductRepo, services.CouponServiceImpl)

	quantity := 2
	request := &requests.PlaceOrderRequest{
		Items: []requests.OrderItemRequest{
			{ProductId: "invalid", Quantity: &quantity},
		},
	}

	result, err := service.PlaceOrder(context.Background(), request)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusBadRequest, err.ErrorCode)
	assert.Contains(t, err.Message, "invalid product id")
}

// TestOrderService_PlaceOrder_MultipleItems tests the PlaceOrder method with multiple items
func TestOrderService_PlaceOrder_MultipleItems(t *testing.T) {
	mockCouponRepo := new(MockCouponRepository)
	mockCouponRepo.On("GetCouponCounts", mock.Anything).Return(int64(100), int64(2), nil)
	mockCouponRepo.On("GetCouponsByFileCount", mock.Anything, ">= 2").Return([]string{"SAVE1000", "DISCOUNT50"}, nil)

	err := services.InitializeCouponService(mockCouponRepo)
	assert.Nil(t, err)

	mockOrderRepo := new(MockOrderRepository)
	mockProductRepo := new(MockProductRepository)
	service := services.NewOrderServiceImpl(mockOrderRepo, mockProductRepo, services.CouponServiceImpl)

	quantity1 := 2
	quantity2 := 3
	request := &requests.PlaceOrderRequest{
		Items: []requests.OrderItemRequest{
			{ProductId: "1", Quantity: &quantity1},
			{ProductId: "2", Quantity: &quantity2},
		},
	}

	mockProducts := []*models.Product{
		{Id: 1, Name: "Product 1", Price: 10.00, Category: "Cat1", Status: "available"},
		{Id: 2, Name: "Product 2", Price: 15.00, Category: "Cat2", Status: "available"},
	}

	mockProductRepo.On("GetByIds", mock.Anything, []int64{1, 2}).Return(mockProducts, nil)
	mockOrderRepo.On("CreateOrder", mock.Anything, mock.AnythingOfType("*models.Order"), mock.AnythingOfType("[]models.OrderItem")).
		Run(func(args mock.Arguments) {
			order := args.Get(1).(*models.Order)
			order.Id = "550e8400-e29b-41d4-a716-446655440002"
		}).
		Return(nil)

	result, err := service.PlaceOrder(context.Background(), request)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.Id)
	assert.Len(t, result.Items, 2)

	mockProductRepo.AssertExpectations(t)
	mockOrderRepo.AssertExpectations(t)
}

// TestOrderService_PlaceOrder_DuplicateItems_ShouldAggregate tests the PlaceOrder method with duplicate items
func TestOrderService_PlaceOrder_DuplicateItems_ShouldAggregate(t *testing.T) {
	mockCouponRepo := new(MockCouponRepository)
	mockCouponRepo.On("GetCouponCounts", mock.Anything).Return(int64(100), int64(2), nil)
	mockCouponRepo.On("GetCouponsByFileCount", mock.Anything, ">= 2").Return([]string{"SAVE1000", "DISCOUNT50"}, nil)

	err := services.InitializeCouponService(mockCouponRepo)
	assert.Nil(t, err)

	mockOrderRepo := new(MockOrderRepository)
	mockProductRepo := new(MockProductRepository)
	service := services.NewOrderServiceImpl(mockOrderRepo, mockProductRepo, services.CouponServiceImpl)

	quantity1 := 2
	quantity2 := 3
	request := &requests.PlaceOrderRequest{
		Items: []requests.OrderItemRequest{
			{ProductId: "1", Quantity: &quantity1},
			{ProductId: "1", Quantity: &quantity2},
		},
	}

	mockProducts := []*models.Product{{Id: 1, Name: "Product 1", Price: 10.00, Category: "Cat1", Status: "available"}}

	mockProductRepo.On("GetByIds", mock.Anything, []int64{1}).Return(mockProducts, nil)
	mockOrderRepo.On("CreateOrder", mock.Anything, mock.AnythingOfType("*models.Order"), mock.AnythingOfType("[]models.OrderItem")).
		Run(func(args mock.Arguments) {
			order := args.Get(1).(*models.Order)
			order.Id = "550e8400-e29b-41d4-a716-446655440003"
		}).
		Return(nil)

	result, err := service.PlaceOrder(context.Background(), request)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.Id)
	assert.Len(t, result.Items, 1)
	assert.Equal(t, 5, result.Items[0].Quantity)

	mockProductRepo.AssertExpectations(t)
	mockOrderRepo.AssertExpectations(t)
}

// TestOrderService_PlaceOrder_CreateOrderFails tests the PlaceOrder method with a failure to create an order
func TestOrderService_PlaceOrder_CreateOrderFails(t *testing.T) {
	mockCouponRepo := new(MockCouponRepository)
	mockCouponRepo.On("GetCouponCounts", mock.Anything).Return(int64(100), int64(2), nil)
	mockCouponRepo.On("GetCouponsByFileCount", mock.Anything, ">= 2").Return([]string{"SAVE1000", "DISCOUNT50"}, nil)

	err := services.InitializeCouponService(mockCouponRepo)
	assert.Nil(t, err)

	mockOrderRepo := new(MockOrderRepository)
	mockProductRepo := new(MockProductRepository)
	service := services.NewOrderServiceImpl(mockOrderRepo, mockProductRepo, services.CouponServiceImpl)

	quantity := 2
	request := &requests.PlaceOrderRequest{
		Items: []requests.OrderItemRequest{
			{ProductId: "1", Quantity: &quantity},
		},
	}

	mockProducts := []*models.Product{{Id: 1, Name: "Product 1", Price: 10.00, Category: "Cat1", Status: "available"}}
	mockError := &errors.ErrorDetails{
		ErrorCode: http.StatusInternalServerError,
		Message:   "failed to fetch products",
	}

	mockProductRepo.On("GetByIds", mock.Anything, []int64{1}).Return(mockProducts, nil)
	mockOrderRepo.On("CreateOrder", mock.Anything, mock.AnythingOfType("*models.Order"), mock.AnythingOfType("[]models.OrderItem")).Return(mockError)

	result, err := service.PlaceOrder(context.Background(), request)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.ErrorCode)

	mockProductRepo.AssertExpectations(t)
	mockOrderRepo.AssertExpectations(t)
}
