#!/bin/bash
set -e

echo "Initializing database..."
cd "$(dirname "$0")/.."
go run cmd/db_init/main.go
echo "Database initialization completed."
