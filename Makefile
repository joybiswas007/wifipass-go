# Define the binary name
BINARY = wifipass

# Define source directories
SRC_DIR = ./cmd/wifipass/

# Compiler flags
GOFLAGS = -ldflags="-s -w"

# Default target
all: build

# Build the binary
build:
	go build $(GOFLAGS) -o $(BINARY) $(SRC_DIR)

# Clean up built files
clean:
	rm -rf $(BINARY)

# Install dependencies
deps:
	go mod tidy

# Format code
fmt:
	go fmt .

# Lint the code (requires golangci-lint)
lint:
	golangci-lint run

# Install golangci-lint if not installed
install-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Display help
help:
	@echo "Available commands:"
	@echo "  make build       - Build the project"
	@echo "  make clean       - Remove built files"
	@echo "  make deps        - Install dependencies"
	@echo "  make fmt         - Format code"
	@echo "  make lint        - Run linter"
	@echo "  make install-lint - Install golangci-lint"

