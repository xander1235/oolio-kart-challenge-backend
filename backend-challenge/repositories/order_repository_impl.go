package repositories

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"net/http"
	"oolio.com/kart/configs"

	"oolio.com/kart/exceptions"
	"oolio.com/kart/exceptions/errors"
	"oolio.com/kart/models"
)

type OrderRepositoryImpl struct {
	pool *pgxpool.Pool
}

// NewOrderRepositoryImpl creates a new instance of OrderRepositoryImpl
func NewOrderRepositoryImpl(pool *pgxpool.Pool) *OrderRepositoryImpl {
	return &OrderRepositoryImpl{pool: pool}
}

// CreateOrder creates a new order in the database
func (o *OrderRepositoryImpl) CreateOrder(ctx context.Context, order *models.Order, items []models.OrderItem) *errors.ErrorDetails {
	tx, err := o.pool.Begin(ctx)
	if err != nil {
		configs.Logger.Error("failed to begin transaction", zap.Error(err))
		return exceptions.GenericException("failed to begin transaction", http.StatusInternalServerError)
	}

	// Marshal order meta
	var metaJSON []byte
	if order.Meta != nil {
		metaJSON, err = json.Marshal(order.Meta)
		if err != nil {
			configs.Logger.Error("failed to marshal order meta", zap.Error(err))
			return exceptions.GenericException("failed to marshal order meta", http.StatusInternalServerError)
		}
	}

	orderQuery := `INSERT INTO orders (coupon_code, subtotal, discount, total, meta)
                   VALUES ($1, $2, $3, $4, $5)
                   RETURNING id, created_at, modified_at`

	err = tx.QueryRow(ctx, orderQuery,
		order.CouponCode,
		order.Subtotal,
		order.Discount,
		order.Total,
		metaJSON,
	).Scan(&order.Id, &order.CreatedAt, &order.ModifiedAt)

	if err != nil {
		txErr := tx.Rollback(ctx)
		if txErr != nil {
			configs.Logger.Error("failed to rollback transaction", zap.Error(txErr))
		}
		configs.Logger.Error("failed to save order", zap.Error(err))
		return exceptions.GenericException("failed to save order", http.StatusInternalServerError)
	}

	if len(items) > 0 {
		itemQuery := `INSERT INTO order_items (order_id, product_id, quantity, unit_price, price, meta)
                      VALUES ($1, $2, $3, $4, $5, $6)
                      RETURNING id, created_at`

		for i := range items {
			var itemMetaJSON []byte
			if items[i].Meta != nil {
				itemMetaJSON, err = json.Marshal(items[i].Meta)
				if err != nil {
					txErr := tx.Rollback(ctx)
					if txErr != nil {
						configs.Logger.Error("failed to rollback transaction", zap.Error(txErr))
					}
					configs.Logger.Error("failed to marshal order item meta", zap.Error(err))
					return exceptions.GenericException("failed to marshal order item meta", http.StatusInternalServerError)
				}
			}

			err = tx.QueryRow(ctx, itemQuery,
				order.Id,
				items[i].ProductId,
				items[i].Quantity,
				items[i].UnitPrice,
				items[i].Price,
				itemMetaJSON,
			).Scan(&items[i].Id, &items[i].CreatedAt)

			if err != nil {
				txErr := tx.Rollback(ctx)
				if txErr != nil {
					configs.Logger.Error("failed to rollback transaction", zap.Error(txErr))
				}
				configs.Logger.Error("failed to save order item", zap.Error(err))
				return exceptions.GenericException("failed to save order item", http.StatusInternalServerError)
			}

			items[i].OrderId = order.Id
		}
	}

	if err = tx.Commit(ctx); err != nil {
		txErr := tx.Rollback(ctx)
		if txErr != nil {
			configs.Logger.Error("failed to rollback transaction", zap.Error(txErr))
		}
		return exceptions.GenericException("failed to commit transaction", http.StatusInternalServerError)
	}

	return nil
}
