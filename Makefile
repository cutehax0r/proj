.PHONY: help build build-macos-arm64 build-macos-amd64 build-windows-arm64 build-windows-amd64 build-linux-arm64 build-linux-amd64 build-all clean test acceptance

# Default target
.DEFAULT_GOAL := build

# Detect local platform
LOCAL_OS := $(shell uname -s | tr '[:upper:]' '[:lower:]' | sed 's/darwin/darwin/;s/linux/linux/')
LOCAL_ARCH := $(shell uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')

# Variables
BINARY_NAME := proj
GO := go
GOOS ?= $(LOCAL_OS)
GOARCH ?= $(LOCAL_ARCH)
OUTPUT_DIR := ./bin
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

# Platform-specific extensions
ifeq ($(GOOS),windows)
BINARY_EXT := .exe
endif

# Help target
help:
	@echo "Makefile targets for $(BINARY_NAME):"
	@echo ""
	@echo "Build targets:"
	@echo "  make build                    Build for local platform ($(LOCAL_OS)/$(LOCAL_ARCH))"
	@echo "  make build-macos-arm64        Build for macOS/arm64"
	@echo "  make build-macos-amd64        Build for macOS/amd64 (Intel)"
	@echo "  make build-windows-arm64      Build for Windows/arm64"
	@echo "  make build-windows-amd64      Build for Windows/amd64 (Intel)"
	@echo "  make build-linux-arm64        Build for Linux/arm64"
	@echo "  make build-linux-amd64        Build for Linux/amd64 (Intel)"
	@echo "  make build-all                Build for all platforms"
	@echo ""
	@echo "Development targets:"
	@echo "  make test                     Run unit tests"
	@echo "  make acceptance               Run acceptance tests (builds + tests)"
	@echo "  make clean                    Remove build artifacts"
	@echo "  make help                     Display this help message"
	@echo ""
	@echo "Environment variables:"
	@echo "  GOOS=<os>                     Target OS (darwin, windows, linux)"
	@echo "  GOARCH=<arch>                 Target architecture (arm64, amd64)"
	@echo ""
	@echo "Examples:"
	@echo "  make build-windows-amd64      Build Windows 64-bit"
	@echo "  make build-all                Build all platform variants"
	@echo "  GOOS=linux GOARCH=arm64 make  Custom build"

# Main build target (builds for local platform + creates copy as 'proj')
build:
	@echo "Building $(BINARY_NAME) for $(GOOS)/$(GOARCH)..."
	@mkdir -p $(OUTPUT_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build -o $(OUTPUT_DIR)/$(BINARY_NAME)-$(GOOS)-$(GOARCH)$(BINARY_EXT) .
	@cp $(OUTPUT_DIR)/$(BINARY_NAME)-$(GOOS)-$(GOARCH)$(BINARY_EXT) $(OUTPUT_DIR)/$(BINARY_NAME)
	@echo "✓ Built: $(OUTPUT_DIR)/$(BINARY_NAME)-$(GOOS)-$(GOARCH)$(BINARY_EXT)"
	@echo "✓ Copied to: $(OUTPUT_DIR)/$(BINARY_NAME)"

# Platform-specific build targets
build-macos-arm64:
	@echo "Building $(BINARY_NAME) for macOS/arm64..."
	@mkdir -p $(OUTPUT_DIR)
	GOOS=darwin GOARCH=arm64 $(GO) build -o $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-arm64 .
	@echo "✓ Built: $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-arm64"

build-macos-amd64:
	@echo "Building $(BINARY_NAME) for macOS/amd64..."
	@mkdir -p $(OUTPUT_DIR)
	GOOS=darwin GOARCH=amd64 $(GO) build -o $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-amd64 .
	@echo "✓ Built: $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-amd64"

build-windows-arm64:
	@echo "Building $(BINARY_NAME) for Windows/arm64..."
	@mkdir -p $(OUTPUT_DIR)
	GOOS=windows GOARCH=arm64 $(GO) build -o $(OUTPUT_DIR)/$(BINARY_NAME)-windows-arm64.exe .
	@echo "✓ Built: $(OUTPUT_DIR)/$(BINARY_NAME)-windows-arm64.exe"

build-windows-amd64:
	@echo "Building $(BINARY_NAME) for Windows/amd64..."
	@mkdir -p $(OUTPUT_DIR)
	GOOS=windows GOARCH=amd64 $(GO) build -o $(OUTPUT_DIR)/$(BINARY_NAME)-windows-amd64.exe .
	@echo "✓ Built: $(OUTPUT_DIR)/$(BINARY_NAME)-windows-amd64.exe"

build-linux-arm64:
	@echo "Building $(BINARY_NAME) for Linux/arm64..."
	@mkdir -p $(OUTPUT_DIR)
	GOOS=linux GOARCH=arm64 $(GO) build -o $(OUTPUT_DIR)/$(BINARY_NAME)-linux-arm64 .
	@echo "✓ Built: $(OUTPUT_DIR)/$(BINARY_NAME)-linux-arm64"

build-linux-amd64:
	@echo "Building $(BINARY_NAME) for Linux/amd64..."
	@mkdir -p $(OUTPUT_DIR)
	GOOS=linux GOARCH=amd64 $(GO) build -o $(OUTPUT_DIR)/$(BINARY_NAME)-linux-amd64 .
	@echo "✓ Built: $(OUTPUT_DIR)/$(BINARY_NAME)-linux-amd64"

# Build all platforms
build-all: build-macos-arm64 build-macos-amd64 build-windows-arm64 build-windows-amd64 build-linux-arm64 build-linux-amd64
	@echo ""
	@echo "✓ Successfully built for all platforms"
	@echo "Binaries in: $(OUTPUT_DIR)/"
	@ls -lh $(OUTPUT_DIR)/

# Test target
test:
	@echo "Running tests..."
	$(GO) test -v ./...

# Acceptance test target
acceptance: build
	@echo "Running acceptance tests..."
	$(GO) test -v -tags=acceptance ./tests/acceptance/...

# Clean target
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(OUTPUT_DIR)
	@$(GO) clean
	@echo "✓ Clean complete"
