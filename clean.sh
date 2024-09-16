#!/bin/bash

# Define paths to clean up the generated files
PROTO_GEN_PATH="proto"

# Find and remove all generated .pb.go and _grpc.pb.go files
echo "Cleaning up generated .pb.go and _grpc.pb.go files..."
find "$PROTO_GEN_PATH" -name "*.pb.go" -type f -delete

# Optionally, remove other temporary or build-related files if needed
# Uncomment the following lines if you want to clean other types of files:
# echo "Cleaning up binaries and build directories..."
# rm -rf bin/ build/

echo "Cleanup completed."
