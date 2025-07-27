include .env
export

DATABASE_URL = "postgresql://$$DB_USERNAME:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=disable"

dev:
	air

build:
	go build -o ./build/main ./cmd/api

start:
	./build/main

migrate-create:
	@read -p "Enter migration name (use underscore): " name; \
	migrate create -ext sql -dir migrations -seq $$name

migrate-up:
	migrate -database $(DATABASE_URL) -path migrations up
	
migrate-down:
	migrate -database $(DATABASE_URL) -path migrations down

migrate-force:
	@read -p "Enter migration version: " version; \
	migrate -database $(DATABASE_URL) -path migrations force $$version

seed:
	go run ./cmd/seeder

seed-factory:
	go run ./cmd/seeder -factory

seed-factory-custom:
	go run ./cmd/seeder -factory -users=50

seed-build:
	go build -o ./build/seeder ./cmd/seeder

seed-run:
	./build/seeder

seed-run-factory:
	./build/seeder -factory