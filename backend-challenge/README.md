# Advanced Challenge

Build an API server implementing our OpenAPI spec for food ordering API in [Go](https://go.dev).\
You can find our [API Documentation](https://orderfoodonline.deno.dev/public/openapi.html) here.

API documentation is based on [OpenAPI3.1](https://swagger.io/specification/v3/) specification.
You can also find spec file [here](https://orderfoodonline.deno.dev/public/openapi.yaml).

> The API immplementation example available to you at orderfoodonline.deno.dev/api is simplified and doesn't handle some edge cases intentionally.
> Use your best judgement to build a Robust API server.

## Basic Requirements

- Implement all APIs described in the OpenAPI specification
- Conform to the OpenAPI specification as close to as possible
- Implement all features our [demo API server](https://orderfoodonline.deno.dev) has implemented
- Validate promo code according to promo code validation logic described below

### Promo Code Validation

You will find three big files containing random text in this repositotory.\
A promo code is valid if the following rules apply:

1. Must be a string of length between 8 and 10 characters
2. It can be found in **at least two** files

> Files containing valid coupons are couponbase1.gz, couponbase2.gz and couponbase3.gz

You can download the files from here

[file 1](https://orderfoodonline-files.s3.ap-southeast-2.amazonaws.com/couponbase1.gz)
[file 2](https://orderfoodonline-files.s3.ap-southeast-2.amazonaws.com/couponbase2.gz)
[file 3](https://orderfoodonline-files.s3.ap-southeast-2.amazonaws.com/couponbase3.gz)

**Example Promo Codes**

Valid promo codes

- HAPPYHRS
- FIFTYOFF

Invalid promo codes

- SUPER100

> [!TIP]
> it should be noted that there are more valid and invalid promo codes that those shown above.

## Getting Started

You might need to configure Git LFS to clone this repository\
https://github.com/oolio-group/kart-challenge/tree/advanced-challenge/backend-challenge

1. Use this repository as a template and create a new repository in your account
2. Start coding
3. Share your repository

---

# Setup & Run Instructions

## Prerequisites

- Go 1.23 or higher
- Docker & Docker Compose
- PostgreSQL 14+ (if running locally without Docker)

## Quick Start with Docker (Recommended)

### Option 1: Using Docker Compose

```bash
docker compose up --build -d
```

This will:
- Build the Go application
- Start PostgreSQL database
- Run database migrations automatically
- Start the API server on port 8080

### Stop the services

```bash
# Using Docker Compose
docker compose down
```

### Access the application
- **API Base URL**: http://localhost:8080/api
- **Health Check**: http://localhost:8080/api/health
- **Swagger UI**: http://localhost:8080/swagger/index.html

---

## Running Locally (Without Docker)

### 1. Install dependencies
```bash
go mod download
```

### 2. Setup PostgreSQL database
```bash

# Run migrations
psql -d postgres -f schemas/schemas.sql
```

### 3. Configure environment variables

Create a `.env` file in the project root:
```env
APP_NAME=kart-api
HOST=localhost
PORT=8080
RELEASE_ENV=development
LOG_LEVEL=info

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=postgres
DB_SCHEMA=kart
DB_SSL_MODE=disable
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME=300
```

### 4. Run the application
```bash
export IS_LOCAL=true && go run main.go
```

---

## Running Tests

### Run all tests
```bash
go test -v ./tests/...
```

### Run with coverage
```bash
go test -v -cover ./tests/...
```

### Run specific test package
```bash
# Controller tests
go test -v ./tests/controllers

# Service tests
go test -v ./tests/services
```

---

## API Endpoints

### Health Check
```bash
curl http://localhost:8080/api/health
```

### Get Products
```bash
curl http://localhost:8080/api/product
```

### Get Product by ID
```bash
curl http://localhost:8080/api/product/1
```

### Place Order
```bash
curl -X POST http://localhost:8080/api/order \
  -H "Content-Type: application/json" \
  -H "api_key: api_test" \
  -d '{
    "items": [
      {
        "productId": "1",
        "quantity": 2
      }
    ]
  }'
```

### Place Order with Coupon
```bash
curl -X POST http://localhost:8080/api/order \
  -H "Content-Type: application/json" \
  -H "api_key: api_test" \
  -d '{
    "couponCode": "HAPPYHRS",
    "items": [
      {
        "productId": "1",
        "quantity": 2
      }
    ]
  }'
```

