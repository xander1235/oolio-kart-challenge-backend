package repositories

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5"
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
	txOptions := pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	}

	tx, err := o.pool.BeginTx(ctx, txOptions)
	if err != nil {
		configs.Logger.Error("failed to begin transaction", zap.Error(err))
		return exceptions.GenericException("failed to begin transaction", http.StatusInternalServerError)
	}

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
		batch := &pgx.Batch{}

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

			batch.Queue(
				`INSERT INTO order_items (order_id, product_id, quantity, unit_price, price, meta)
				 VALUES ($1, $2, $3, $4, $5, $6)
				 RETURNING id, created_at`,
				order.Id,
				items[i].ProductId,
				items[i].Quantity,
				items[i].UnitPrice,
				items[i].Price,
				itemMetaJSON,
			)
		}

		batchResults := tx.SendBatch(ctx, batch)
		defer batchResults.Close()

		for i := range items {
			err = batchResults.QueryRow().Scan(&items[i].Id, &items[i].CreatedAt)
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
		return exceptions.GenericException("failed to commit transaction", http.StatusInternalServerError)
	}

	return nil
}
