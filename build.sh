#!/bin/bash

echo ">> Install swag CLI"
go install github.com/swaggo/swag/cmd/swag@latest

echo ">> Generate Swagger docs..."
swag init

echo ">> Build Go binary..."
go build -o out main.go
