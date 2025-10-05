package responses

import (
	"oolio.com/kart/models"
	"strconv"
)

// ProductResponse represents a product in the API response
type ProductResponse struct {
	Id       string  `json:"id" example:"1" doc:"Unique product ID"`
	Name     string  `json:"name" example:"Margherita Pizza" doc:"Product name"`
	Category string  `json:"category" example:"Pizza" doc:"Product category"`
	Price    float64 `json:"price" example:"12.99" doc:"Product price in USD"`
} //@name Product

// ToProductResponse converts domain model to API response
func ToProductResponse(product *models.Product) *ProductResponse {
	return &ProductResponse{
		Id:       strconv.Itoa(int(product.Id)),
		Name:     product.Name,
		Price:    product.Price,
		Category: product.Category,
	}
}

// ToProductResponses converts multiple domain models to API responses
func ToProductResponses(products []*models.Product) []*ProductResponse {
	responses := make([]*ProductResponse, len(products))
	for i, product := range products {
		responses[i] = ToProductResponse(product)
	}
	return responses
}
