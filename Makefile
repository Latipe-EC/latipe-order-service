API_NAME = order_service
WORKER_NAME = worker_name
GO= go
BUILD_DIR = \build

API_MAIN_FILE = ./cmd/main.go
API_SERVER_FILE = ./internal/server.go

WORKER_MAIN_FILE = ./cmd/main.go
WORKER_SERVER_FILE = ./internal/server.go

setup:
	go mod tidy
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/swaggo/swag/cmd/swag@latest

wire-worker:
	cd order-worker/internal/ && wire
wire-api:
	cd order-rest-api/internal/ && wire

# build binary
build-app:
	make wire-worker
	make wire-api
	cd order-worker/ && $(GO) build -o $(BUILD_DIR)\ $(API_MAIN_FILE)
	cd order-rest-api/ && $(GO) build -o $(BUILD_DIR)\ $(WORKER_MAIN_FILE)

run-worker:
	make wire-worker
	cd order-worker/ && $(GO) run $(API_MAIN_FILE)

run-api:
	make wire-api
	cd order-rest-api/ && $(GO) run $(WORKER_MAIN_FILE)