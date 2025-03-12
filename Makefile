PROTO_SRC_DIR = ./protobuf
PROTO_OUT_DIR = ./model 

.PHONY: prepare-proto-gen
prepare-proto-gen:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

.PHONY: proto-gen
proto-gen:
	protoc --proto_path=$(PROTO_SRC_DIR) \
	    --go_out=$(PROTO_OUT_DIR) \
	    --go_opt=paths=source_relative \
		--go-grpc_out=. \
	    $(PROTO_SRC_DIR)/block.proto
