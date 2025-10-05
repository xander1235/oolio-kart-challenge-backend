package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Product struct {
	ID       int64                  `json:"id"`
	Name     string                 `json:"name"`
	Category string                 `json:"category"`
	Price    float64                `json:"price"`
	Status   string                 `json:"status"`
	Image    map[string]interface{} `json:"image"`
	Meta     map[string]interface{} `json:"meta,omitempty"`
}

type ProductMigration struct {
	pool         *pgxpool.Pool
	dataFilePath string
}

func NewProductMigration(pool *pgxpool.Pool, dataFilePath string) *ProductMigration {
	return &ProductMigration{
		pool:         pool,
		dataFilePath: dataFilePath,
	}
}

func (pm *ProductMigration) Run(ctx context.Context) error {
	needsMigration, err := pm.checkMigrationNeeded(ctx)
	if err != nil {
		return fmt.Errorf("failed to check migration status: %w", err)
	}

	if !needsMigration {
		log.Println("Products already loaded, skipping migration")
		return nil
	}

	log.Println("Starting product migration")

	products, err := pm.loadProductsFromJSON()
	if err != nil {
		return fmt.Errorf("failed to load products from JSON: %w", err)
	}

	if err := pm.insertProducts(ctx, products); err != nil {
		return fmt.Errorf("failed to insert products: %w", err)
	}

	log.Printf("Product migration completed: %d products loaded", len(products))
	return nil
}

func (pm *ProductMigration) checkMigrationNeeded(ctx context.Context) (bool, error) {
	var count int64
	err := pm.pool.QueryRow(ctx, "SELECT COUNT(*) FROM products").Scan(&count)
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

func (pm *ProductMigration) loadProductsFromJSON() ([]Product, error) {
	data, err := os.ReadFile(pm.dataFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read product file: %w", err)
	}

	var products []Product
	if err := json.Unmarshal(data, &products); err != nil {
		return nil, fmt.Errorf("failed to unmarshal products: %w", err)
	}

	return products, nil
}

func (pm *ProductMigration) insertProducts(ctx context.Context, products []Product) error {
	for _, product := range products {
		imageJSON, err := json.Marshal(product.Image)
		if err != nil {
			return fmt.Errorf("failed to marshal image for product %s: %w", product.Name, err)
		}

		var metaJSON []byte
		if product.Meta != nil {
			metaJSON, err = json.Marshal(product.Meta)
			if err != nil {
				return fmt.Errorf("failed to marshal meta for product %s: %w", product.Name, err)
			}
		}

		query := `
			INSERT INTO products (name, category, price, status, image, meta)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (name, category) DO UPDATE
			SET price = EXCLUDED.price,
			    status = EXCLUDED.status,
			    image = EXCLUDED.image,
			    meta = EXCLUDED.meta,
			    modified_at = NOW()
		`

		_, err = pm.pool.Exec(ctx, query,
			product.Name,
			product.Category,
			product.Price,
			"available",
			imageJSON,
			metaJSON,
		)

		if err != nil {
			return fmt.Errorf("failed to insert product %s: %w", product.Name, err)
		}
	}

	return nil
}
