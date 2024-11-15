# Binary name
BINARY_NAME=tx-parser

# Go related variables
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin

# Main package path
MAIN_PACKAGE=cmd/server

# Build variables
BUILD_DIR=${GOBASE}/build
BUILD_PACKAGE=github.com/vdhieu/tx-parser/${MAIN_PACKAGE}

.PHONY: all build clean test coverage lint run

.PHONY: all build clean test coverage lint run mock

all: clean lint test mock build


mock:
	@echo "Generating mocks..."
	@if command -v mockery >/dev/null; then \
		mockery --all --keeptree --output mocks --inpackage; \
	else \
		echo "mockery is not installed. Installing mockery..."; \
		go install github.com/vektra/mockery/v2@latest; \
		mockery --all --keeptree --output mocks --inpackage; \
	fi

build:
	@echo "Building..."
	@go build -o ${BUILD_DIR}/${BINARY_NAME} ${BUILD_PACKAGE}


clean:
	@echo "Cleaning..."
	@rm -rf ${BUILD_DIR}
	@go clean

test:
	@echo "Running tests..."
	@go test -v ./...

coverage:
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out

lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint is not installed. Please install it first."; \
		exit 1; \
	fi

run:
	@echo "Running application..."
	@go run ${BUILD_PACKAGE}
