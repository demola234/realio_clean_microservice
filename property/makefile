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



PROTO_DIR = interfaces/api/proto
GRPC_DIR = interfaces/api/grpc
PB_DIR = interfaces/api/grpc
OUT_DIR = infrastructure/api/grpc
SWAGGER_DIR = docs/swagger

proto:
	@rm -f $(GRPC_DIR)/*.go
	@rm -f docs/swagger/*.swagger.json
	protoc --proto_path=infrastructure/api/proto \
	       --proto_path=infrastructure/api/proto/google \
	       --go_out=infrastructure/api/grpc --go_opt=paths=source_relative \
	       --go-grpc_out=infrastructure/api/grpc --go-grpc_opt=paths=source_relative \
	       --grpc-gateway_out=infrastructure/api/grpc --grpc-gateway_opt=paths=source_relative \
	       --openapiv2_out=docs/swagger --openapiv2_opt=allow_merge,merge_file_name=realio-property \
	       infrastructure/api/proto/*.proto
	go run scripts/fix_swagger.go
	statik -src=./docs/swagger -dest=./docs

.PHONY: proto


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
	docker run --name $(DB_NAME) -p $(DB_PORT):5444 -e POSTGRES_USER=$(DB_USER) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -d postgres:15-alpine

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

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir db/migration -seq $${name}

sqlc:
	sqlc generate

server:
	go run cmd/property/main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/demola234/property/db/sqlc Store

evans:
	evans --host localhost --port 9090 -r repl