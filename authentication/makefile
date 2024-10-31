PROTO_DIR = interfaces/api/proto
GRPC_DIR = interfaces/api/grpc
PB_DIR = interfaces/api/grpc
DB_URL = "postgresql://usmiekrnx6pqq6ybbp2g:IeEtTWSo6FGjweshdN9yUco06Kieba@bjhuxkyprn3x6g7vhtnl-postgresql.services.clever-cloud.com:50013/bjhuxkyprn3x6g7vhtnl"

proto:
	@rm -f $(GRPC_DIR)/*.go
	protoc --proto_path=$(PROTO_DIR) \
	       --proto_path=$(PROTO_DIR)/google \
	       --go_out=$(GRPC_DIR) --go_opt=paths=source_relative \
	       --go-grpc_out=$(GRPC_DIR) --go-grpc_opt=paths=source_relative \
	       --grpc-gateway_out=$(GRPC_DIR) --grpc-gateway_opt=paths=source_relative \
	       $(PROTO_DIR)/*.proto

.PHONY: proto

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