#!/bin/bash

# Navigate to the cmd/app where main.go is located
cd ../cmd/app || exit

# Run the swag init command
swag init --parseDependencyLevel 3 --output ../../pkg/api/v1/docs
