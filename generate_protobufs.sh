#!/bin/bash

# Navigate to the src directory
cd src

# Create the output directory if it doesn't exist
mkdir -p pb

# Find all .proto files and generate Go code
find ./protos -name "*.proto" -type f -exec protoc \
    --go_out=. \
    {} \; || {
    echo "Error: Failed to generate protobuf files"
    exit 1
}

echo "Successfully generated protobuf files"
