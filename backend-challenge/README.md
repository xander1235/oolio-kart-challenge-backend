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

## Quick Start with Docker Compose (Recommended)

### 1. Clone the repository
```bash
git clone <your-repo-url>
cd backend-challenge
```

### 2. Start the application
```bash
docker-compose up --build
```

This will:
- Build the Go application
- Start PostgreSQL database
- Run database migrations automatically
- Start the API server on port 8080

### 3. Access the application
- **API Base URL**: http://localhost:8080/api
- **Health Check**: http://localhost:8080/api/health
- **Swagger UI**: http://localhost:8080/swagger/index.html

### 4. Stop the application
```bash
docker-compose down
```

To remove all data and start fresh:
```bash
docker-compose down -v
```

---

## Running Locally (Without Docker)

### 1. Install dependencies
```bash
go mod download
```

### 2. Setup PostgreSQL database
```bash
# Create database
createdb kart_db

# Run migrations
psql -d kart_db -f schemas/schemas.sql
```

### 3. Configure environment variables

Create a `.env` file in the project root:
```env
IS_LOCAL=true
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
DB_NAME=kart_db
DB_SSL_MODE=disable
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME=300
```

### 4. Run the application
```bash
go run main.go
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

---

## Project Structure

```
backend-challenge/
├── configs/           # Application configuration
├── constants/         # Application constants
├── controllers/       # HTTP request handlers
├── dtos/             # Data transfer objects
├── exceptions/       # Error handling
├── middlewares/      # HTTP middlewares
├── models/           # Domain models
├── repositories/     # Database layer
├── routes/           # Route definitions
├── services/         # Business logic
├── schemas/          # Database schemas
├── tests/            # Unit tests
├── docs/             # Swagger documentation
├── Dockerfile        # Docker image definition
├── docker-compose.yml # Docker compose configuration
└── main.go           # Application entry point
```

---

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `IS_LOCAL` | Load from .env file | `false` |
| `APP_NAME` | Application name | `kart-api` |
| `HOST` | Server host | `localhost` |
| `PORT` | Server port | `8080` |
| `RELEASE_ENV` | Environment (development/production) | `development` |
| `LOG_LEVEL` | Log level (debug/info/warn/error) | `info` |
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port | `5432` |
| `DB_USER` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | - |
| `DB_NAME` | Database name | `kart_db` |
| `DB_SSL_MODE` | SSL mode | `disable` |

---

## Troubleshooting

### Port already in use
```bash
# Find and kill process using port 8080
lsof -ti:8080 | xargs kill -9
```

### Database connection issues
```bash
# Check PostgreSQL is running
docker-compose ps

# View logs
docker-compose logs postgres
```

### Reset database
```bash
docker-compose down -v
docker-compose up --build
```

---

## Development

### Generate Swagger docs
```bash
swag init
```

### Format code
```bash
go fmt ./...
```

### Run linter
```bash
golangci-lint run
```

---

## Notes

- The coupon validation feature requires large data files (not included in Git)
- Without coupon data, all coupon codes will be rejected (expected behavior)
- API key middleware is enabled for order endpoints (use `api_key: api_test`)
- Database migrations run automatically on first startup with Docker Compose
