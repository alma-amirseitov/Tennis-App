.PHONY: dev build test test-coverage lint migrate-up migrate-down migrate-create sqlc seed

include .env
export

# Development
dev:
	air

build:
	go build -o bin/server cmd/server/main.go

# Testing
test:
	go test ./... -v -count=1

test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

lint:
	golangci-lint run ./...

# Database
migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down 1

migrate-create:
	migrate create -ext sql -dir migrations -seq $(NAME)

# Code generation
sqlc:
	sqlc generate

# Data
seed:
	go run scripts/seed/main.go
