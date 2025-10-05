package controllers_test

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"oolio.com/kart/controllers"
	"oolio.com/kart/dtos/requests"
	"oolio.com/kart/dtos/responses"
	"oolio.com/kart/exceptions/errors"
	"testing"
)

// MockOrderService is a mock implementation of OrderService
func TestOrderController_PlaceOrder_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockOrderService)
	controller := controllers.NewOrderController(mockService)

	quantity := 2
	requestBody := requests.PlaceOrderRequest{
		Items: []requests.OrderItemRequest{
			{ProductId: "1", Quantity: &quantity},
		},
	}

	mockResponse := &responses.OrderResponse{
		Id: "550e8400-e29b-41d4-a716-446655440000",
		Items: []responses.OrderItemResponse{
			{ProductId: "1", Quantity: 2},
		},
	}

	mockService.On("PlaceOrder", mock.Anything, mock.AnythingOfType("*requests.PlaceOrderRequest")).Return(mockResponse, nil)

	router := gin.New()
	router.POST("/orders", controller.PlaceOrder)

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response responses.OrderResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", response.Id)

	mockService.AssertExpectations(t)
}

// MockOrderService is a mock implementation of OrderService
func TestOrderController_PlaceOrder_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockOrderService)
	controller := controllers.NewOrderController(mockService)

	router := gin.New()
	router.POST("/orders", controller.PlaceOrder)

	req, _ := http.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "invalid_request", response["type"])
}

// MockOrderService is a mock implementation of OrderService
func TestOrderController_PlaceOrder_EmptyItems(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockOrderService)
	controller := controllers.NewOrderController(mockService)

	requestBody := requests.PlaceOrderRequest{
		Items: []requests.OrderItemRequest{},
	}

	router := gin.New()
	router.POST("/orders", controller.PlaceOrder)

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// MockOrderService is a mock implementation of OrderService
func TestOrderController_PlaceOrder_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockOrderService)
	controller := controllers.NewOrderController(mockService)

	quantity := 2
	requestBody := requests.PlaceOrderRequest{
		Items: []requests.OrderItemRequest{
			{ProductId: "999", Quantity: &quantity},
		},
	}

	mockError := &errors.ErrorDetails{
		ErrorCode: http.StatusBadRequest,
		Message:   "product not found",
	}

	mockService.On("PlaceOrder", mock.Anything, mock.AnythingOfType("*requests.PlaceOrderRequest")).Return(nil, mockError)

	router := gin.New()
	router.POST("/orders", controller.PlaceOrder)

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "product not found", response["message"])

	mockService.AssertExpectations(t)
}

// MockOrderService is a mock implementation of OrderService
func TestOrderController_PlaceOrder_WithCoupon(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockOrderService)
	controller := controllers.NewOrderController(mockService)

	quantity := 2
	requestBody := requests.PlaceOrderRequest{
		CouponCode: "SAVE10",
		Items: []requests.OrderItemRequest{
			{ProductId: "1", Quantity: &quantity},
		},
	}

	mockResponse := &responses.OrderResponse{
		Id:         "550e8400-e29b-41d4-a716-446655440000",
		CouponCode: "SAVE10",
		Items: []responses.OrderItemResponse{
			{ProductId: "1", Quantity: 2},
		},
	}

	mockService.On("PlaceOrder", mock.Anything, mock.AnythingOfType("*requests.PlaceOrderRequest")).Return(mockResponse, nil)

	router := gin.New()
	router.POST("/orders", controller.PlaceOrder)

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response responses.OrderResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "SAVE10", response.CouponCode)

	mockService.AssertExpectations(t)
}
