#!/bin/bash
set -e
cd golang
go mod tidy
go run ./cmd/main.go
