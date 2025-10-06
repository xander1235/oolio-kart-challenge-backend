package services

import (
	"context"
	"go.uber.org/zap"
	"net/http"
	"oolio.com/kart/configs"
	"oolio.com/kart/exceptions"
	"strconv"

	"oolio.com/kart/dtos/requests"
	"oolio.com/kart/dtos/responses"
	"oolio.com/kart/exceptions/errors"
	"oolio.com/kart/models"
	repoBase "oolio.com/kart/repositories/base"
	serviceBase "oolio.com/kart/services/base"
)

type OrderServiceImpl struct {
	orderRepository       repoBase.OrderRepository
	productRepository     repoBase.ProductRepository
	couponService         serviceBase.CouponService
	maxQuantityPerProduct int
}

// NewOrderServiceImpl creates a new instance of OrderServiceImpl
func NewOrderServiceImpl(orderRepository repoBase.OrderRepository, productRepository repoBase.ProductRepository, couponService serviceBase.CouponService) *OrderServiceImpl {
	return &OrderServiceImpl{
		orderRepository:       orderRepository,
		productRepository:     productRepository,
		couponService:         couponService,
		maxQuantityPerProduct: 1000,
	}
}

// PlaceOrder places a new order
func (s *OrderServiceImpl) PlaceOrder(ctx context.Context, request *requests.PlaceOrderRequest) (*responses.OrderResponse, *errors.ErrorDetails) {
	if request.CouponCode != "" {
		if s.couponService == nil {
			configs.Logger.Error("coupon service not available")
			return nil, exceptions.GenericException("some internal error occurred", http.StatusInternalServerError)
		}

		isValid, err := s.couponService.ValidateCoupon(ctx, request.CouponCode)
		if err != nil {
			configs.Logger.Error("coupon service validate coupon code", zap.Any("error", err))
			return nil, err
		}

		if !isValid {
			configs.Logger.Error("invalid coupon code")
			return nil, exceptions.UnprocessableEntityException("invalid coupon code")
		}
	}

	itemMap := make(map[string]*models.OrderItem)
	for _, reqItem := range request.Items {
		if existing, found := itemMap[reqItem.ProductId]; found {
			newQty := existing.Quantity + *reqItem.Quantity
			if newQty > s.maxQuantityPerProduct {
				configs.Logger.Error("quantity exceeds maximum limit")
				return nil, exceptions.BadRequestException("quantity exceeds maximum limit")
			}
			existing.Quantity = newQty
		} else {
			productId, err := strconv.ParseInt(reqItem.ProductId, 10, 64)
			if err != nil {
				configs.Logger.Error("invalid product id", zap.Any("error", err))
				return nil, exceptions.BadRequestException("invalid product id")
			}

			if *reqItem.Quantity > s.maxQuantityPerProduct {
				configs.Logger.Error("quantity exceeds maximum limit", zap.Any("error", err))
				return nil, exceptions.BadRequestException("quantity exceeds maximum limit")
			}

			domainItem := models.OrderItem{
				ProductId: productId,
				Quantity:  *reqItem.Quantity,
			}
			itemMap[reqItem.ProductId] = &domainItem
		}
	}

	var aggregatedItems []models.OrderItem
	for _, item := range itemMap {
		aggregatedItems = append(aggregatedItems, *item)
	}

	var subtotal float64
	var discount float64 = 0
	var products []*models.Product

	productIds := make([]int64, 0, len(aggregatedItems))
	for i := range aggregatedItems {
		productIds = append(productIds, aggregatedItems[i].ProductId)
	}

	const batchSize = 100
	for i := 0; i < len(productIds); i += batchSize {
		end := i + batchSize
		if end > len(productIds) {
			end = len(productIds)
		}
		chunk := productIds[i:end]
		chunkProducts, err := s.productRepository.GetByIds(ctx, chunk)
		if err != nil {
			configs.Logger.Error("products not found", zap.Any("error", err))
			return nil, exceptions.BadRequestException("products not found")
		}
		products = append(products, chunkProducts...)
	}

	productMap := make(map[int64]*models.Product)
	for _, product := range products {
		productMap[product.Id] = product
	}

	for i := range aggregatedItems {
		product, exists := productMap[aggregatedItems[i].ProductId]
		if !exists {
			return nil, exceptions.BadRequestException("product not found")
		}

		aggregatedItems[i].UnitPrice = product.Price
		aggregatedItems[i].Price = product.Price * float64(aggregatedItems[i].Quantity)
		subtotal += aggregatedItems[i].Price
	}

	// TODO Discount logic comes here
	var total = subtotal
	if discount > 0 {
		// TODO Apply discount with limit
		total -= discount
	}

	order := &models.Order{
		CouponCode: request.CouponCode,
		Subtotal:   subtotal,
		Discount:   discount,
		Total:      total,
	}

	if saveErr := s.orderRepository.CreateOrder(ctx, order, aggregatedItems); saveErr != nil {
		configs.Logger.Error("failed to create order", zap.Any("error", saveErr))
		return nil, saveErr
	}

	response := responses.ToOrderResponse(order, aggregatedItems, products)
	return response, nil
}
