package requests

// PlaceOrderRequest represents the request to place an order
type PlaceOrderRequest struct {
	CouponCode string             `json:"couponCode,omitempty" example:"HAPPYHRS" doc:"Optional coupon code for discount"`
	Items      []OrderItemRequest `json:"items" binding:"required,min=1,dive" doc:"List of items to order (minimum 1 item required)"`
} //@name OrderReq

// OrderItemRequest represents an item in the order request
type OrderItemRequest struct {
	ProductId string `json:"productId" binding:"required" example:"1" doc:"Product ID to order"`
	Quantity  *int   `json:"quantity" binding:"required,min=1" example:"2" doc:"Quantity to order (minimum 1)"`
} //@name OrderItemReq
