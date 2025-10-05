#!/bin/bash

set -e

echo "==================================="
echo "Kart API - Database Migration Tool"
echo "==================================="
echo ""

# Check if .env exists
if [ ! -f .env ]; then
    echo "⚠️  .env file not found!"
    echo "Creating .env from .env.example..."
    cp .env.example .env
    echo "✅ Created .env file"
    echo "⚠️  Please edit .env with your database credentials"
    echo ""
    exit 1
fi

# Parse arguments
MIGRATION_TYPE="${1:-all}"

echo "Migration type: $MIGRATION_TYPE"
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed"
    echo "Please install Go 1.23 or higher"
    exit 1
fi

echo "✅ Go version: $(go version)"
echo ""

# Download dependencies
echo "📦 Downloading dependencies..."
go mod download
echo "✅ Dependencies downloaded"
echo ""

# Run migration
echo "🚀 Starting migration..."
echo ""

go run . -type="$MIGRATION_TYPE"

EXIT_CODE=$?

echo ""
if [ $EXIT_CODE -eq 0 ]; then
    echo "✅ Migration completed successfully!"
else
    echo "❌ Migration failed with exit code: $EXIT_CODE"
    exit $EXIT_CODE
fi
