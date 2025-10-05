# Database Migrations

Standalone migration scripts for Kart API database.

## Prerequisites

- Go 1.23+
- PostgreSQL 14+
- Coupon data files (optional, can download from S3)

## Setup

### 1. Install dependencies

```bash
cd migrations
go mod download
```

### 2. Configure environment

Copy `.env.example` to `.env` and update values:

```bash
cp .env.example .env
```

Edit `.env`:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=kart_db
DB_SSL_MODE=disable

# For local files
COUPON_SOURCE=local
COUPON_DATA_DIR=../data

# For S3 download
# COUPON_SOURCE=s3
# COUPON_S3_BASE_URL=https://orderfoodonline-files.s3.ap-southeast-2.amazonaws.com

COUPON_FORCE_MIGRATION=false
PRODUCT_DATA_FILE=../data/product.json
```

## Running Migrations

### Run all migrations

```bash
go run . -type=all
```

### Run specific migrations

```bash
# Product migration only
go run . -type=product

# Coupon migration only
go run . -type=coupon
```

### Using custom .env file

```bash
go run . -type=all -env=.env.production
```

### Build and run

```bash
# Build binary
go build -o migrate

# Run
./migrate -type=all
```

## Migration Types

### Product Migration

- Loads products from `product.json`
- Inserts into `products` table
- Skips if products already exist
- Fast (< 1 second)

### Coupon Migration

- Loads coupons from 3 gzip files (local or S3)
- Processes ~100M+ coupons
- Uses parallel processing (16 partitions, 8 workers)
- Creates staging table for performance
- Aggregates and deduplicates
- Duration: 10-30 minutes depending on hardware

## Data Sources

### Local Files

Place files in `../data/` directory:
- `couponbase1.gz`
- `couponbase2.gz`
- `couponbase3.gz`
- `product.json`

### S3 Download

Set `COUPON_SOURCE=s3` in `.env`. Files will be streamed from S3 without downloading.