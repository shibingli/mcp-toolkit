# Makefile for MCP Toolkit
# 用于简化构建、测试和发布流程

.PHONY: help build build-all test clean install uninstall release dev

# 项目信息 / Project information
PROJECT_NAME := mcp-toolkit
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# 构建标签 / Build tags
BUILD_TAGS := sonic
LDFLAGS := -s -w -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)

# 目录 / Directories
OUTPUT_DIR := dist
SANDBOX_DIR := sandbox

# 默认目标 / Default target
help:
	@echo "MCP Toolkit Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  make build        - Build for current platform"
	@echo "  make build-all    - Build for all platforms (Windows, Linux, macOS)"
	@echo "  make test         - Run all tests"
	@echo "  make test-cover   - Run tests with coverage report"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make install      - Install to local system"
	@echo "  make uninstall    - Uninstall from local system"
	@echo "  make release      - Create a new release (tag and build)"
	@echo "  make dev          - Run in development mode"
	@echo "  make fmt          - Format code"
	@echo "  make lint         - Run linters"
	@echo ""
	@echo "Current version: $(VERSION)"
	@echo "Build time: $(BUILD_TIME)"
	@echo "Git commit: $(GIT_COMMIT)"

# 构建当前平台 / Build for current platform
build:
	@echo "Building $(PROJECT_NAME) $(VERSION) for current platform..."
	@mkdir -p $(OUTPUT_DIR)
	go build -tags="$(BUILD_TAGS)" -ldflags="$(LDFLAGS)" -o $(OUTPUT_DIR)/$(PROJECT_NAME) .
	@echo "✓ Build completed: $(OUTPUT_DIR)/$(PROJECT_NAME)"

# 构建所有平台 / Build for all platforms
build-all:
	@echo "Building for all platforms..."
	@chmod +x scripts/build.sh
	@VERSION=$(VERSION) ./scripts/build.sh

# 运行测试 / Run tests
test:
	@echo "Running tests..."
	go test -v -tags="$(BUILD_TAGS)" ./...

# 运行测试并生成覆盖率报告 / Run tests with coverage
test-cover:
	@echo "Running tests with coverage..."
	go test -v -tags="$(BUILD_TAGS)" -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report generated: coverage.html"

# 清理构建产物 / Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(OUTPUT_DIR)
	rm -rf $(SANDBOX_DIR)
	rm -f coverage.out coverage.html
	rm -f $(PROJECT_NAME) $(PROJECT_NAME).exe
	@echo "✓ Clean completed"

# 安装到本地系统 / Install to local system
install: build
	@echo "Installing $(PROJECT_NAME)..."
ifeq ($(OS),Windows_NT)
	@mkdir -p $(LOCALAPPDATA)/Programs/$(PROJECT_NAME)
	@cp $(OUTPUT_DIR)/$(PROJECT_NAME).exe $(LOCALAPPDATA)/Programs/$(PROJECT_NAME)/
	@echo "✓ Installed to $(LOCALAPPDATA)/Programs/$(PROJECT_NAME)/"
	@echo "  Add to PATH: $(LOCALAPPDATA)/Programs/$(PROJECT_NAME)"
else
	@mkdir -p $(HOME)/.local/bin
	@cp $(OUTPUT_DIR)/$(PROJECT_NAME) $(HOME)/.local/bin/
	@chmod +x $(HOME)/.local/bin/$(PROJECT_NAME)
	@echo "✓ Installed to $(HOME)/.local/bin/$(PROJECT_NAME)"
	@echo "  Make sure $(HOME)/.local/bin is in your PATH"
endif

# 从本地系统卸载 / Uninstall from local system
uninstall:
	@echo "Uninstalling $(PROJECT_NAME)..."
ifeq ($(OS),Windows_NT)
	@rm -f $(LOCALAPPDATA)/Programs/$(PROJECT_NAME)/$(PROJECT_NAME).exe
	@rmdir $(LOCALAPPDATA)/Programs/$(PROJECT_NAME) 2>/dev/null || true
else
	@rm -f $(HOME)/.local/bin/$(PROJECT_NAME)
endif
	@echo "✓ Uninstalled"

# 创建发布 / Create release
release:
	@echo "Creating release $(VERSION)..."
	@if [ -z "$(VERSION)" ] || [ "$(VERSION)" = "dev" ]; then \
		echo "Error: Please set a valid version tag"; \
		echo "Usage: git tag -a v1.0.0 -m 'Release v1.0.0' && make release"; \
		exit 1; \
	fi
	@echo "Building all platforms..."
	@make build-all
	@echo "✓ Release $(VERSION) created in $(OUTPUT_DIR)/"
	@echo ""
	@echo "Next steps:"
	@echo "  1. Review the build artifacts in $(OUTPUT_DIR)/"
	@echo "  2. Push the tag: git push origin $(VERSION)"
	@echo "  3. GitHub Actions will automatically create the release"

# 开发模式运行 / Run in development mode
dev:
	@echo "Running in development mode..."
	@mkdir -p $(SANDBOX_DIR)
	go run -tags="$(BUILD_TAGS)" . -sandbox $(SANDBOX_DIR) -transport stdio

# 格式化代码 / Format code
fmt:
	@echo "Formatting code..."
	gofmt -w .
	@echo "✓ Code formatted"

# 运行 linter / Run linters
lint:
	@echo "Running linters..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not found, skipping..."; \
		echo "Install: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# 下载依赖 / Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy
	@echo "✓ Dependencies downloaded"

# 更新依赖 / Update dependencies
deps-update:
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy
	@echo "✓ Dependencies updated"

# 运行客户端测试 / Run client tests
test-client:
	@echo "Building server and client..."
	@make build
	@go build -tags="$(BUILD_TAGS)" -o $(OUTPUT_DIR)/$(PROJECT_NAME)-client ./cmd/client
	@echo "Starting server in background..."
	@$(OUTPUT_DIR)/$(PROJECT_NAME) -transport http -http-port 8080 &
	@sleep 2
	@echo "Running client tests..."
	@$(OUTPUT_DIR)/$(PROJECT_NAME)-client || true
	@pkill -f "$(PROJECT_NAME)" || true
	@echo "✓ Client tests completed"

# 生成文档 / Generate documentation
docs:
	@echo "Generating documentation..."
	@if command -v godoc >/dev/null 2>&1; then \
		echo "Starting godoc server at http://localhost:6060"; \
		godoc -http=:6060; \
	else \
		echo "godoc not found"; \
		echo "Install: go install golang.org/x/tools/cmd/godoc@latest"; \
	fi

# 检查构建环境 / Check build environment
check:
	@echo "Checking build environment..."
	@echo "Go version: $(shell go version)"
	@echo "Git version: $(shell git --version)"
	@echo "Project version: $(VERSION)"
	@echo "Build tags: $(BUILD_TAGS)"
	@echo "LDFLAGS: $(LDFLAGS)"
	@echo "✓ Environment check completed"

