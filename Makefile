APP_NAME=ffaas
CMD_DIR=./cmd/ffaas-server
BIN_DIR=bin

.PHONY: run build test fmt lint clean

## Run the server (memory backend by default)
run:
	@echo ">> Running $(APP_NAME)..."
	go run $(CMD_DIR)

## Build binary
build:
	@echo ">> Building binary..."
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(APP_NAME) $(CMD_DIR)

## Run tests
test:
	@echo ">> Running tests..."
	go test ./... -v

## Format code
fmt:
	@echo ">> Formatting code..."
	go fmt ./...

## Lint (basic vet)
lint:
	@echo ">> Vetting code..."
	go vet ./...

## Clean build artifacts
clean:
	@echo ">> Cleaning..."
	rm -rf $(BIN_DIR)

