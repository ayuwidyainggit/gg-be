# Build And Development
include .env

dev:
	@go run main.go

doc:
	@swag init

dev-reload:
	@air

install:
	@go mod tidy

migration:
	migrate create -ext sql -dir migrations $(command)

migrateUp:
	migrate -path migrations -database $(DATABASE_URL) -verbose up

migrateDown:
	migrate -path migrations -database $(DATABASE_URL) -verbose down

migrateForce:
	migrate -path migrations -database $(DATABASE_URL) -verbose force $(command)

migrateDrop:
	migrate -path migrations -database $(DATABASE_URL) -verbose drop

.PHONY: dev doc dev-reload install migration migrateUp migrateDown migrateForce migrateDrop
