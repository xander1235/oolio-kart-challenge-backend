package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"net/http"
	"oolio.com/kart/configs"

	"oolio.com/kart/exceptions"
	"oolio.com/kart/exceptions/errors"
	"oolio.com/kart/models"
)

type ProductRepositoryImpl struct {
	pool *pgxpool.Pool
}

func NewProductRepositoryImpl(pool *pgxpool.Pool) *ProductRepositoryImpl {
	return &ProductRepositoryImpl{pool: pool}
}

// Save Saves a new product to the database
func (p *ProductRepositoryImpl) Save(ctx context.Context, product *models.Product) *errors.ErrorDetails {
	imageJSON, err := json.Marshal(product.Image)
	if err != nil {
		configs.Logger.Error("failed to marshal product image", zap.Error(err))
		return exceptions.GenericException("failed to marshal product image", http.StatusInternalServerError)
	}

	var metaJSON []byte
	if product.Meta != nil {
		metaJSON, err = json.Marshal(product.Meta)
		if err != nil {
			configs.Logger.Error("failed to marshal product meta", zap.Error(err))
			return exceptions.GenericException("failed to marshal product meta", http.StatusInternalServerError)
		}
	}

	query := `INSERT INTO products (name, category, price, status, image, meta)
              VALUES ($1, $2, $3, $4, $5, $6)
              RETURNING id, created_at, modified_at`

	err = p.pool.QueryRow(ctx, query,
		product.Name,
		product.Category,
		product.Price,
		product.Status,
		imageJSON,
		metaJSON,
	).Scan(&product.Id, &product.CreatedAt, &product.ModifiedAt)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			configs.Logger.Error("product already exists", zap.Error(err))
			return exceptions.BadRequestException("product already exists")
		}
		configs.Logger.Error("failed to save product", zap.Error(err))
		return exceptions.GenericException("failed to save product", http.StatusInternalServerError)
	}

	return nil
}

// Update Updates an existing product in the database
func (p *ProductRepositoryImpl) Update(ctx context.Context, product *models.Product) *errors.ErrorDetails {
	if product.Id == 0 {
		configs.Logger.Error("product id is required")
		return exceptions.BadRequestException("product id is required")
	}

	imageJSON, err := json.Marshal(product.Image)
	if err != nil {
		configs.Logger.Error("failed to marshal product image", zap.Error(err))
		return exceptions.GenericException("failed to marshal product image", http.StatusInternalServerError)
	}

	var metaJSON []byte
	if product.Meta != nil {
		metaJSON, err = json.Marshal(product.Meta)
		if err != nil {
			configs.Logger.Error("failed to marshal product meta", zap.Error(err))
			return exceptions.GenericException("failed to marshal product meta", http.StatusInternalServerError)
		}
	}

	query := `UPDATE products
              SET name = $1,
                  category = $2,
                  price = $3,
                  status = $4,
                  image = $5,
                  meta = $6,
                  modified_at = NOW()
              WHERE id = $7
              RETURNING created_at, modified_at`

	err = p.pool.QueryRow(ctx, query,
		product.Name,
		product.Category,
		product.Price,
		product.Status,
		imageJSON,
		metaJSON,
		product.Id,
	).Scan(&product.CreatedAt, &product.ModifiedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			configs.Logger.Error("product not found", zap.Error(err))
			return exceptions.GenericException("product not found", http.StatusNotFound)
		}
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			configs.Logger.Error("product with this name and category already exists", zap.Error(err))
			return exceptions.GenericException("product with this name and category already exists", http.StatusConflict)
		}
		configs.Logger.Error("failed to update product", zap.Error(err))
		return exceptions.GenericException("failed to update product", http.StatusInternalServerError)
	}

	return nil
}

// GetById Retrieves a product by its ID from the database
func (p *ProductRepositoryImpl) GetById(ctx context.Context, id int64) (*models.Product, *errors.ErrorDetails) {
	query := `SELECT id, name, category, price, status, image, meta, created_at, modified_at
              FROM products
              WHERE id = $1`

	row := p.pool.QueryRow(ctx, query, id)
	product, err := scanProduct(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			configs.Logger.Error("product not found", zap.Error(err))
			return nil, exceptions.GenericException("product not found", http.StatusNotFound)
		}
		configs.Logger.Error("failed to fetch product", zap.Error(err))
		return nil, exceptions.GenericException("failed to fetch product", http.StatusInternalServerError)
	}

	return product, nil
}

// ListProducts Retrieves a list of products from the database
func (p *ProductRepositoryImpl) ListProducts(ctx context.Context, offset, limit *int) []*models.Product {
	query := `SELECT id, name, category, price, status, image, meta, created_at, modified_at
              FROM products
              ORDER BY id`

	var args []any
	paramIdx := 1

	if offset != nil {
		query += fmt.Sprintf(" OFFSET $%d", paramIdx)
		args = append(args, *offset)
		paramIdx++
	}

	if limit != nil {
		query += fmt.Sprintf(" LIMIT $%d", paramIdx)
		args = append(args, *limit)
	}

	rows, err := p.pool.Query(ctx, query, args...)
	if err != nil {
		return []*models.Product{}
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		product, scanErr := scanProduct(rows)
		if scanErr != nil {
			continue
		}
		products = append(products, product)
	}

	if rows.Err() != nil {
		return []*models.Product{}
	}

	return products
}

// GetByIds Retrieves a list of products by their IDs from the database
func (p *ProductRepositoryImpl) GetByIds(ctx context.Context, ids []int64) ([]*models.Product, *errors.ErrorDetails) {
	query := `SELECT id, name, category, price, status, image, meta, created_at, modified_at
              FROM products
              WHERE id = ANY($1)`

	rows, err := p.pool.Query(ctx, query, ids)
	if err != nil {
		return []*models.Product{}, exceptions.GenericException("failed to fetch products", http.StatusInternalServerError)
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		product, scanErr := scanProduct(rows)
		if scanErr != nil {
			continue
		}
		products = append(products, product)
	}

	if rows.Err() != nil {
		return []*models.Product{}, exceptions.GenericException("failed to fetch products", http.StatusInternalServerError)
	}

	return products, nil
}

func scanProduct(row pgx.Row) (*models.Product, error) {
	var (
		imageBytes []byte
		metaBytes  []byte
	)

	product := &models.Product{}

	err := row.Scan(
		&product.Id,
		&product.Name,
		&product.Category,
		&product.Price,
		&product.Status,
		&imageBytes,
		&metaBytes,
		&product.CreatedAt,
		&product.ModifiedAt,
	)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(imageBytes, &product.Image); err != nil {
		return nil, err
	}

	if len(metaBytes) > 0 {
		var meta map[string]any
		if err = json.Unmarshal(metaBytes, &meta); err != nil {
			return nil, err
		}
		product.Meta = meta
	}

	return product, nil
}
