# Stage 1: Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code (data folder excluded via .dockerignore)
COPY . .

# Build the application
# CGO_ENABLED=0 for static binary
# -ldflags="-s -w" to strip debug info and reduce size
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" \
    -o kart-api \
    .

# Stage 2: Runtime stage
FROM alpine:latest

# Install runtime dependencies (including wget for healthcheck)
RUN apk --no-cache add ca-certificates tzdata wget

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/kart-api .

# Copy schemas folder for reference (optional, if needed at runtime)
COPY --from=builder /build/schemas ./schemas

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget -q -O - http://localhost:8080/api/health || exit 1

# Run the application
CMD ["./kart-api"]