package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"oolio.com/kart/dtos/responses"
	"oolio.com/kart/services/base"
	"strconv"
)

type ProductController struct {
	productService base.ProductService
}

// NewProductController creates a new instance of ProductController
func NewProductController(productService base.ProductService) *ProductController {
	return &ProductController{
		productService: productService,
	}
}

// GetProducts godoc
// @Summary      Get all products
// @Description  Retrieve a list of all products
// @Tags         products
// @Produce      json
// @Success      200 {array} responses.ProductResponse
// @Router       /product [get]
func (p *ProductController) GetProducts(c *gin.Context) {
	limitQueryStr := c.Query("limit")
	offsetQueryStr := c.Query("offset")

	var limit, offset *int
	if limitQueryStr != "" {
		parsedLimit, err := strconv.Atoi(limitQueryStr)
		if err == nil {
			limit = &parsedLimit
		}
	}

	if offsetQueryStr != "" {
		parsedOffset, err := strconv.Atoi(offsetQueryStr)
		if err == nil {
			offset = &parsedOffset
		}
	}

	products := p.productService.GetProducts(c.Request.Context(), limit, offset)
	c.JSON(http.StatusOK, responses.ToProductResponses(products))
}

// GetProductById godoc
// @Summary      Get product by ID
// @Description  Retrieve a single product by its ID
// @Tags         products
// @Produce      json
// @Param        productId path int true "Product ID"
// @Success      200 {object} responses.ProductResponse
// @Failure      400 {object} responses.APIResponse
// @Failure      404 {object} responses.APIResponse
// @Failure      500 {object} responses.APIResponse
// @Router       /product/{productId} [get]
func (p *ProductController) GetProductById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("productId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Code:    http.StatusBadRequest,
			Type:    "validation_error",
			Message: "invalid product id",
		})
		return
	}

	product, internalErr := p.productService.GetProductById(c.Request.Context(), id)
	if internalErr != nil {
		c.JSON(internalErr.ErrorCode, responses.APIResponse{
			Code:    internalErr.ErrorCode,
			Type:    "error",
			Message: internalErr.Message,
		})
		return
	}
	c.JSON(http.StatusOK, responses.ToProductResponse(product))
}
