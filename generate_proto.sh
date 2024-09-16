#!/bin/bash

# Command to generate Go and gRPC code from the proto file
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/cacophony.proto

echo "Protobuf and gRPC Go code generation completed."
