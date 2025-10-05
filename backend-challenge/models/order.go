package models

import (
	"time"
)

// Order represents a customer order
type Order struct {
	Id         string         `json:"id"`
	CouponCode string         `json:"coupon_code,omitempty"`
	Subtotal   float64        `json:"subtotal,omitempty"`
	Discount   float64        `json:"discount,omitempty"`
	Total      float64        `json:"total,omitempty"`
	Meta       map[string]any `json:"meta,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	ModifiedAt time.Time      `json:"modified_at"`
}

// OrderItem represents a line item in an order
type OrderItem struct {
	Id         int64          `json:"id,omitempty"`
	OrderId    string         `json:"order_id,omitempty"`
	ProductId  int64          `json:"product_id"`
	Quantity   int            `json:"quantity"`
	UnitPrice  float64        `json:"unit_price,omitempty"`
	Price      float64        `json:"price,omitempty"`
	Meta       map[string]any `json:"meta,omitempty"`
	CreatedAt  time.Time      `json:"created_at,omitempty"`
	ModifiedAt time.Time      `json:"modified_at"`
}
