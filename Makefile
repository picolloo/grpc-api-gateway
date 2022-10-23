.PHONY: proto

proto:
	protoc -I ./proto/ \
	--go_out=./proto \
	--go-grpc_out=./proto \
	./proto/request.proto
