package controllers_test

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"oolio.com/kart/controllers"
	"oolio.com/kart/exceptions/errors"
	"oolio.com/kart/models"
	"testing"
)

// MockProductService is a mock implementation of ProductService
func TestProductController_GetProducts_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockProductService)
	controller := controllers.NewProductController(mockService)

	mockProducts := []*models.Product{
		{
			Id:       1,
			Name:     "Margherita Pizza",
			Price:    12.99,
			Category: "Pizza",
			Status:   "available",
		},
		{
			Id:       2,
			Name:     "Pepperoni Pizza",
			Price:    14.99,
			Category: "Pizza",
			Status:   "available",
		},
	}

	mockService.On("GetProducts", mock.Anything, mock.Anything, mock.Anything).Return(mockProducts)

	router := gin.New()
	router.GET("/products", controller.GetProducts)

	req, _ := http.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, "Margherita Pizza", response[0]["name"])

	mockService.AssertExpectations(t)
}

// MockProductService is a mock implementation of ProductService
func TestProductController_GetProducts_WithPagination(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockProductService)
	controller := controllers.NewProductController(mockService)

	mockProducts := []*models.Product{
		{Id: 1, Name: "Product 1", Price: 10.00, Category: "Category1", Status: "available"},
	}

	mockService.On("GetProducts", mock.Anything, mock.Anything, mock.Anything).Return(mockProducts)

	router := gin.New()
	router.GET("/products", controller.GetProducts)

	req, _ := http.NewRequest(http.MethodGet, "/products?limit=10&offset=0", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

// MockProductService is a mock implementation of ProductService
func TestProductController_GetProductById_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockProductService)
	controller := controllers.NewProductController(mockService)

	mockProduct := &models.Product{
		Id:       1,
		Name:     "Margherita Pizza",
		Price:    12.99,
		Category: "Pizza",
		Status:   "available",
	}

	mockService.On("GetProductById", mock.Anything, int64(1)).Return(mockProduct, nil)

	router := gin.New()
	router.GET("/products/:productId", controller.GetProductById)

	req, _ := http.NewRequest(http.MethodGet, "/products/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Margherita Pizza", response["name"])
	assert.Equal(t, 12.99, response["price"])

	mockService.AssertExpectations(t)
}

// MockProductService is a mock implementation of ProductService
func TestProductController_GetProductById_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockProductService)
	controller := controllers.NewProductController(mockService)

	router := gin.New()
	router.GET("/products/:productId", controller.GetProductById)

	req, _ := http.NewRequest(http.MethodGet, "/products/invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "validation_error", response["type"])
	assert.Equal(t, "invalid product id", response["message"])
}

// MockProductService is a mock implementation of ProductService
func TestProductController_GetProductById_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockProductService)
	controller := controllers.NewProductController(mockService)

	mockError := &errors.ErrorDetails{
		ErrorCode: http.StatusNotFound,
		Message:   "product not found",
	}

	mockService.On("GetProductById", mock.Anything, int64(999)).Return(nil, mockError)

	router := gin.New()
	router.GET("/products/:productId", controller.GetProductById)

	req, _ := http.NewRequest(http.MethodGet, "/products/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "product not found", response["message"])

	mockService.AssertExpectations(t)
}
