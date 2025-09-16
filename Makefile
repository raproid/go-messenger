# Secure Messenger Makefile

# Variables
APP_NAME = secure-messenger
SERVER_BINARY = messenger-server
CLIENT_BINARY = messenger-client
GO_VERSION = 1.21
BUILD_DIR = bin
DATA_DIR = data
CERT_DIR = certs

# Go build flags
LDFLAGS = -ldflags "-s -w"
BUILD_FLAGS = -trimpath $(LDFLAGS)

# Default target
.PHONY: all
all: build

# Build all binaries
.PHONY: build
build: build-server build-client

# Build server
.PHONY: build-server
build-server:
	@echo "🔨 Building server..."
	@mkdir -p $(BUILD_DIR)
	@cd server && CGO_ENABLED=1 go build $(BUILD_FLAGS) -o ../$(BUILD_DIR)/$(SERVER_BINARY) .
	@echo "✅ Server built: $(BUILD_DIR)/$(SERVER_BINARY)"


# Build client
.PHONY: build-client
build-client:
	@echo "🔨 Building client..."
	@mkdir -p $(BUILD_DIR)
	@cd client_gui && go build $(BUILD_FLAGS) -o ../$(BUILD_DIR)/$(CLIENT_BINARY) .
	@echo "✅ Client built: $(BUILD_DIR)/$(CLIENT_BINARY)"

# Build for multiple platforms
.PHONY: build-all
build-all: build-linux build-windows build-darwin

# Build for Linux
.PHONY: build-linux
build-linux:
	@echo "🔨 Building for Linux..."
	@mkdir -p $(BUILD_DIR)/linux
	@cd server && CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o ../$(BUILD_DIR)/linux/$(SERVER_BINARY) .
	@cd client_gui && GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o ../$(BUILD_DIR)/linux/$(CLIENT_BINARY) .
	@echo "✅ Linux binaries built in $(BUILD_DIR)/linux/"

# Build for Windows
.PHONY: build-windows
build-windows:
	@echo "🔨 Building for Windows..."
	@mkdir -p $(BUILD_DIR)/windows
	@cd server && CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build $(BUILD_FLAGS) -o ../$(BUILD_DIR)/windows/$(SERVER_BINARY).exe .
	@cd client_gui && GOOS=windows GOARCH=amd64 go build $(BUILD_FLAGS) -o ../$(BUILD_DIR)/windows/$(CLIENT_BINARY).exe .
	@echo "✅ Windows binaries built in $(BUILD_DIR)/windows/"

# Build for macOS
.PHONY: build-darwin
build-darwin:
	@echo "🔨 Building for macOS..."
	@mkdir -p $(BUILD_DIR)/darwin
	@cd server && CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) -o ../$(BUILD_DIR)/darwin/$(SERVER_BINARY) .
	@cd client_gui && GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) -o ../$(BUILD_DIR)/darwin/$(CLIENT_BINARY) .
	@echo "✅ macOS binaries built in $(BUILD_DIR)/darwin/"

# Install dependencies
.PHONY: deps
deps:
	@echo "�� Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "✅ Dependencies installed"

# Generate certificates
.PHONY: certs
certs:
	@echo "🔐 Generating TLS certificates..."
	@mkdir -p $(CERT_DIR)
	@openssl req -x509 -newkey rsa:4096 -keyout $(CERT_DIR)/server.key -out $(CERT_DIR)/server.crt -days 365 -nodes \
		-subj "/C=US/ST=State/L=City/O=SecureMessenger/CN=localhost"
	@echo "✅ Certificates generated in $(CERT_DIR)/"

# Create data directory
.PHONY: data-dir
data-dir:
	@echo "�� Creating data directory..."
	@mkdir -p $(DATA_DIR)
	@echo "✅ Data directory created: $(DATA_DIR)/"

# Run server
.PHONY: run-server
run-server: build-server data-dir
	@echo "🚀 Starting server..."
	@cd $(BUILD_DIR) && ./$(SERVER_BINARY)


# Run client
.PHONY: run-client
run-client: build-client
	@echo "🚀 Starting client..."
	@cd $(BUILD_DIR) && ./$(CLIENT_BINARY)

# Run both server and client
.PHONY: run
run: run-server run-client

# Test
.PHONY: test
test:
	@echo "🧪 Running tests..."
	@go test -v ./...
	@echo "✅ Tests completed"

# Test with coverage
.PHONY: test-coverage
test-coverage:
	@echo "🧪 Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

# Lint
.PHONY: lint
lint:
	@echo "�� Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint not installed, skipping linting"; \
	fi
	@echo "✅ Linting completed"

# Format code
.PHONY: fmt
fmt:
	@echo "🎨 Formatting code..."
	@go fmt ./...
	@echo "✅ Code formatted"

# Clean build artifacts
.PHONY: clean
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -rf $(DATA_DIR)
	@rm -rf $(CERT_DIR)
	@rm -f coverage.out coverage.html
	@echo "✅ Clean completed"

# Clean dependencies
.PHONY: clean-deps
clean-deps:
	@echo "🧹 Cleaning dependencies..."
	@go clean -modcache
	@echo "✅ Dependencies cleaned"

# Install tools
.PHONY: install-tools
install-tools:
	@echo "��️  Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	@echo "✅ Tools installed"

# Security scan
.PHONY: security
security:
	@echo "�� Running security scan..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "⚠️  gosec not installed, run 'make install-tools' first"; \
	fi
	@echo "✅ Security scan completed"

# Docker build
.PHONY: docker-build
docker-build:
	@echo "🐳 Building Docker image..."
	@docker build -t $(APP_NAME):latest .
	@echo "✅ Docker image built: $(APP_NAME):latest"

# Docker run
.PHONY: docker-run
docker-run: docker-build
	@echo "�� Running Docker container..."
	@docker run -p 8080:8080 $(APP_NAME):latest

# Development setup
.PHONY: dev-setup
dev-setup: deps certs data-dir
	@echo "🛠️  Development setup completed"
	@echo "📋 Next steps:"
	@echo "   1. Run 'make run-server' to start the server"
	@echo "   2. Run 'make run-client' to start the client"
	@echo "   3. Or run 'make run' to start both"

# Production setup
.PHONY: prod-setup
prod-setup: deps certs data-dir build-all
	@echo "�� Production setup completed"
	@echo "📋 Binaries available in $(BUILD_DIR)/"

# Help
.PHONY: help
help:
	@echo "Secure Messenger - Available Commands:"
	@echo ""
	@echo "Build Commands:"
	@echo "  build          Build all binaries"
	@echo "  build-server   Build server only"
	@echo "  build-client   Build client only"
	@echo "  build-all      Build for all platforms (Linux, Windows, macOS)"
	@echo ""
	@echo "Run Commands:"
	@echo "  run-server     Run server"
	@echo "  run-client     Run client"
	@echo "  run            Run both server and client"
	@echo ""
	@echo "Development Commands:"
	@echo "  deps           Install dependencies"
	@echo "  test           Run tests"
	@echo "  test-coverage  Run tests with coverage report"
	@echo "  lint           Run linter"
	@echo "  fmt            Format code"
	@echo "  security       Run security scan"
	@echo ""
	@echo "Setup Commands:"
	@echo "  certs          Generate TLS certificates"
	@echo "  data-dir       Create data directory"
	@echo "  dev-setup      Complete development setup"
	@echo "  prod-setup     Complete production setup"
	@echo ""
	@echo "Clean Commands:"
	@echo "  clean          Clean build artifacts"
	@echo "  clean-deps     Clean dependencies"
	@echo ""
	@echo "Docker Commands:"
	@echo "  docker-build   Build Docker image"
	@echo "  docker-run     Run Docker container"
	@echo ""
	@echo "Other Commands:"
	@echo "  install-tools  Install development tools"
	@echo "  help           Show this help message"

# Default target
.DEFAULT_GOAL := help
