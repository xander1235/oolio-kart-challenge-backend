package services_test

import (
	"context"
	"net/http"
	"oolio.com/kart/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"oolio.com/kart/exceptions/errors"
	"oolio.com/kart/models"
)

// TestProductService_GetProducts_Success tests the GetProducts method of the ProductService
func TestProductService_GetProducts_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := services.NewProductServiceImpl(mockRepo)

	mockProducts := []*models.Product{
		{Id: 1, Name: "Product 1", Price: 10.00, Category: "Category1", Status: "available"},
		{Id: 2, Name: "Product 2", Price: 20.00, Category: "Category2", Status: "available"},
	}

	mockRepo.On("ListProducts", mock.Anything, mock.Anything, mock.Anything).Return(mockProducts)

	result := service.GetProducts(context.Background(), nil, nil)

	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, "Product 1", result[0].Name)
	assert.Equal(t, "Product 2", result[1].Name)

	mockRepo.AssertExpectations(t)
}

// TestProductService_GetProducts_WithPagination tests the GetProducts method of the ProductService with pagination
func TestProductService_GetProducts_WithPagination(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := services.NewProductServiceImpl(mockRepo)

	limit := 10
	offset := 0
	mockProducts := []*models.Product{
		{Id: 1, Name: "Product 1", Price: 10.00, Category: "Category1", Status: "available"},
	}

	mockRepo.On("ListProducts", mock.Anything, &offset, &limit).Return(mockProducts)

	result := service.GetProducts(context.Background(), &offset, &limit)

	assert.NotNil(t, result)
	assert.Len(t, result, 1)

	mockRepo.AssertExpectations(t)
}

// TestProductService_GetProducts_EmptyResult tests the GetProducts method of the ProductService with an empty result
func TestProductService_GetProducts_EmptyResult(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := services.NewProductServiceImpl(mockRepo)

	mockRepo.On("ListProducts", mock.Anything, mock.Anything, mock.Anything).Return([]*models.Product{})

	result := service.GetProducts(context.Background(), nil, nil)

	assert.NotNil(t, result)
	assert.Len(t, result, 0)

	mockRepo.AssertExpectations(t)
}

// TestProductService_GetProductById_Success tests the GetProductById method of the ProductService
func TestProductService_GetProductById_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := services.NewProductServiceImpl(mockRepo)

	mockProduct := &models.Product{
		Id:       1,
		Name:     "Margherita Pizza",
		Price:    12.99,
		Category: "Pizza",
		Status:   "available",
	}

	mockRepo.On("GetById", mock.Anything, int64(1)).Return(mockProduct, nil)

	result, err := service.GetProductById(context.Background(), 1)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(1), result.Id)
	assert.Equal(t, "Margherita Pizza", result.Name)
	assert.Equal(t, 12.99, result.Price)

	mockRepo.AssertExpectations(t)
}

// TestProductService_GetProductById_NotFound tests the GetProductById method of the ProductService
func TestProductService_GetProductById_NotFound(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := services.NewProductServiceImpl(mockRepo)

	mockError := &errors.ErrorDetails{
		ErrorCode: http.StatusNotFound,
		Message:   "product not found",
	}

	mockRepo.On("GetById", mock.Anything, int64(999)).Return(nil, mockError)

	result, err := service.GetProductById(context.Background(), 999)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusNotFound, err.ErrorCode)
	assert.Equal(t, "product not found", err.Message)

	mockRepo.AssertExpectations(t)
}
