#!/bin/bash

set -e

PROTO_DIR="proto"
GEN_DIR="gen"

mkdir -p "$GEN_DIR"

find "$PROTO_DIR" -name "*.proto" | while read -r proto_file; do
  echo "generating $proto_file"
  protoc \
    --go_out="$GEN_DIR" \
    --go_opt=paths=source_relative \
    --go-grpc_out="$GEN_DIR" \
    --go-grpc_opt=paths=source_relative \
    --proto_path="." \
    "$proto_file"
done

echo "proto generation complete"
