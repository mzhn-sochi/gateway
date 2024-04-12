BUILD_DIR = ./bin

build:
	mkdir $(BUILD_DIR) #
	go mod tidy
	go build -o ./bin -v ./cmd/gateway

gen:
	protoc -I ./proto ./proto/*.proto --go_out=. --go-grpc_out=.

.DEFAULT_GOAL := run