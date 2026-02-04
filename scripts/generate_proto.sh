#!/bin/bash
set -e

# Install plugins if not present
if ! command -v protoc-gen-go &> /dev/null; then
    echo "Installing protoc-gen-go..."
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
fi
if ! command -v protoc-gen-go-grpc &> /dev/null; then
    echo "Installing protoc-gen-go-grpc..."
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
fi

# Ensure output directories exist
mkdir -p api/generated/common
mkdir -p api/generated/benchmark
mkdir -p api/generated/transaction
mkdir -p api/generated/node
mkdir -p api/generated/metric

# Generate
protoc --proto_path=. \
    --go_out=. --go_opt=module=github.com/fffeng99999/hcp-server \
    --go-grpc_out=. --go-grpc_opt=module=github.com/fffeng99999/hcp-server \
    api/proto/*.proto

echo "Protobuf generation complete."
