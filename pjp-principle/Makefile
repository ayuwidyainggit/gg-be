# Build And Development
include development.env

dev:
	@go run main.go

doc:
	@swag init

dev-reload:
	@air

install:
	@go mod tidy

migration:
	migrate create -ext sql -dir database/migrate -seq $(name)

migrateUp:
	migrate -path database/migrate -database $(DATABASE_URL) -verbose up

migrateDown:
	migrate -path database/migrate -database $(DATABASE_URL) -verbose down

migrateForce:
	migrate -path database/migrate -database $(DATABASE_URL) -verbose force $(version)

migrateDrop:
	migrate -path database/migrate -database $(DATABASE_URL) -verbose drop

.PHONY: dev doc dev-reload install migration migrateUp migrateDown migrateForce migrateDrop
