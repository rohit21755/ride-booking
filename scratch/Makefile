PROTO_DIR := ../proto
PROTO_SRC := ../proto/*.proto
GO_OUT := .

generate-proto:
	protoc \
		--proto_path=$$(PROTO_DIR) \
		--go_out=$$(GO_OUT) \
		--go_opt=module=ride-sharing \
		--go-grpc_out=$$(GO_OUT) \
		--go-grpc_opt=module=ride-sharing \
		$$(PROTO_SRC)
