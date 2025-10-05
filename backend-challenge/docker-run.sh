#!/bin/bash

# Build the Docker image
echo "Building Docker image..."
docker build -t kart-api:latest .

if [ $? -ne 0 ]; then
    echo "Docker build failed!"
    exit 1
fi

echo "Docker image built successfully!"

# Remove existing container
if docker ps -a --format "{{.Names}}" | grep -q "kart-api"; then
    docker rm -f kart-api
fi

# Run the container
echo "Starting kart-api container..."
docker run -d -p 8080:8080 \
  -e IS_LOCAL=false \
  -e APP_NAME=kart-api \
  -e HOST=localhost \
  -e PORT=8080 \
  -e RELEASE_ENV=development \
  -e LOG_LEVEL=info \
  -e DB_HOST=host.docker.internal \
  -e DB_PORT=5432 \
  -e DB_NAME=postgres \
  -e DB_SSL_MODE=disable \
  -e DB_MAX_OPEN_CONNS=25 \
  -e DB_MAX_IDLE_CONNS=5 \
  -e DB_CONN_MAX_LIFETIME=300 \
  --name kart-api \
  kart-api:latest

if [ $? -ne 0 ]; then
    echo "Failed to start container!"
    exit 1
fi

echo "Container started successfully!"
echo "API is running at http://localhost:8080"
echo "Health check: http://localhost:8080/api/health"
echo "Swagger UI: http://localhost:8080/swagger/index.html"
echo ""
echo "To view logs: docker logs -f kart-api"
echo "To stop: docker stop kart-api"
echo "To remove: docker rm kart-api"
