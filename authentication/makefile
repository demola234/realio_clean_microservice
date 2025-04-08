PROTO_DIR = interfaces/api/proto
GRPC_DIR = interfaces/api/grpc
PB_DIR = interfaces/api/grpc

proto:
	@rm -f $(GRPC_DIR)/*.go
	protoc --proto_path=$(PROTO_DIR) \
	       --proto_path=$(PROTO_DIR)/google \
	       --go_out=$(GRPC_DIR) --go_opt=paths=source_relative \
	       --go-grpc_out=$(GRPC_DIR) --go-grpc_opt=paths=source_relative \
	       --grpc-gateway_out=$(GRPC_DIR) --grpc-gateway_opt=paths=source_relative \
	       $(PROTO_DIR)/*.proto

.PHONY: proto

# Load .env file if it exists
ifneq (,$(wildcard .env))
    include .env
    export
endif

DB_URL = ${DB_SOURCE}
DB_NAME ?= ${DB_NAME}
DB_USER ?= ${DB_USER}
DB_PASSWORD ?= ${DB_PASSWORD}
DB_PORT ?= 5433


# Docker commands
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

docker-ps:
	docker-compose ps

docker-build:
	docker-compose build

docker-restart:
	docker-compose restart

# Database commands
postgres:
	docker run --name $(DB_NAME) -p $(DB_PORT):5432 -e POSTGRES_USER=$(DB_USER) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -d postgres:15-alpine

createdb:
	docker exec -it $(DB_NAME) createdb --username=$(DB_USER) --owner=$(DB_USER) $(DB_NAME)

migrateup:
	migrate -database $(DB_URL) -path db/migration -verbose up

migrateupf:
	migrate -database $(DB_URL) -path db/migration force 1

migrateup1:
	migrate -database $(DB_URL) -path db/migration -verbose up 1

migratedown:
	migrate -database $(DB_URL) -path db/migration -verbose down

migratedown1:
	migrate -database $(DB_URL) -path db/migration -verbose down 1

migrateversion:
	migrate -database "$(DB_URL)" -path db/migration version

sqlc:
	sqlc generate

server:
	go run cmd/authentication/main.go