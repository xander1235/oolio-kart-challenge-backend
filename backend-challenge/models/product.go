package models

import "time"

// Product represents a food item available for order
type Product struct {
	Id         int64          `json:"id"`
	Name       string         `json:"name"`
	Image      Image          `json:"image"`
	Price      float64        `json:"price"`
	Category   string         `json:"category"`
	Status     string         `json:"status"`
	Meta       map[string]any `json:"meta,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	ModifiedAt time.Time      `json:"modified_at"`
}

// Image represents the image metadata of a product
type Image struct {
	Thumbnail string `json:"thumbnail"`
	Mobile    string `json:"mobile"`
	Desktop   string `json:"desktop"`
	Tablet    string `json:"tablet"`
}
