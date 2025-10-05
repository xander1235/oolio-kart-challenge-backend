package responses

import (
	"oolio.com/kart/models"
	"strconv"
)

// OrderResponse represents the response after placing an order
type OrderResponse struct {
	CouponCode string              `json:"couponCode" example:"SAVE1000" doc:"Coupon code used for the order"`
	Items      []OrderItemResponse `json:"items" doc:"List of items in the order"`
	Id         string              `json:"id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"Unique order ID (UUID)"`
	Products   []*ProductResponse  `json:"products" doc:"Detailed product information for each item"`
} //@name Order

// OrderItemResponse represents a line item in the order response
type OrderItemResponse struct {
	ProductId string `json:"productId" example:"1" doc:"Product ID"`
	Quantity  int    `json:"quantity" example:"2" doc:"Quantity ordered"`
} //@name OrderItem

// ToOrderResponse converts domain models to API response
func ToOrderResponse(order *models.Order, items []models.OrderItem, products []*models.Product) *OrderResponse {
	itemResponses := make([]OrderItemResponse, len(items))
	for i, item := range items {
		itemResponses[i] = OrderItemResponse{
			ProductId: strconv.Itoa(int(item.ProductId)),
			Quantity:  item.Quantity,
		}
	}

	return &OrderResponse{
		Id:         order.Id,
		Items:      itemResponses,
		Products:   ToProductResponses(products),
		CouponCode: order.CouponCode,
	}
}
