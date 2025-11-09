#!/bin/bash
set -e

echo "Starting local environment with Docker..."
cd "$(dirname "$0")/deployment/local"
docker compose up -d

echo "Running Go app..."
cd ../../cmd/api
go run main.go