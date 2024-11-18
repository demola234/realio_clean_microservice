PROTO_DIR = infrastructure/api/proto
GRPC_DIR = infrastructure/api/grpc
PB_DIR = infrastructure/api/grpc

proto:
	@rm -f $(GRPC_DIR)/*.go
	protoc --proto_path=$(PROTO_DIR) \
	       --proto_path=$(PROTO_DIR)/google \
	       --go_out=$(GRPC_DIR) --go_opt=paths=source_relative \
	       --go-grpc_out=$(GRPC_DIR) --go-grpc_opt=paths=source_relative \
	       --grpc-gateway_out=$(GRPC_DIR) --grpc-gateway_opt=paths=source_relative \
	       $(PROTO_DIR)/*.proto

.PHONY: proto
