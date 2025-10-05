package controllers

import (
	"net/http"
	"oolio.com/kart/dtos/responses"

	"github.com/gin-gonic/gin"

	"oolio.com/kart/dtos/requests"
	"oolio.com/kart/services/base"
)

type OrderController struct {
	orderService base.OrderService
}

// NewOrderController creates a new order controller
func NewOrderController(orderService base.OrderService) *OrderController {
	return &OrderController{orderService: orderService}
}

// PlaceOrder handles POST /api/order
// @Summary      Place a new order
// @Description  Create a new order with items and optional coupon code
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        request body requests.PlaceOrderRequest true "Order details"
// @Success      200 {object} responses.OrderResponse
// @Failure      400 {object} responses.APIResponse
// @Failure      500 {object} responses.APIResponse
// @Security     ApiKeyAuth
// @Param        api_key	  header    string    true   	"api_key must be set for authentication"
// @Router       /order [post]
func (oc *OrderController) PlaceOrder(c *gin.Context) {
	var request requests.PlaceOrderRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Code:    http.StatusBadRequest,
			Type:    "invalid_request",
			Message: err.Error(),
		})
		return
	}

	response, errDetails := oc.orderService.PlaceOrder(c.Request.Context(), &request)
	if errDetails != nil {
		c.JSON(errDetails.ErrorCode, responses.APIResponse{
			Code:    errDetails.ErrorCode,
			Type:    "error",
			Message: errDetails.Message,
		})
		return
	}

	c.JSON(http.StatusOK, response)
}
