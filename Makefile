BINARY_NAME=farmers_backend
MAIN_PATH=./main.go

.PHONY: help build run clean test deps fmt vet

help:
	@echo "  make build    - Build the application"
	@echo "  make run      - Run the application"
	@echo "  make test     - Run tests"
	@echo "  make clean    - Remove build artifacts"
	@echo "  make deps     - Download dependencies"
	@echo "  make fmt      - Format code"
	@echo "  make vet      - Run go vet"
	@echo "  make all      - Format, vet, and build"

build:
	go build -o $(BINARY_NAME) $(MAIN_PATH)

run:
	go run ${MAIN_PATH}

test:
	go test -v ./...

clean:
	@go clean
	@rm -f $(BINARY_NAME)

deps:
	go mod download
	go mod tidy

fmt:
	go fmt ./...

vet:
	go vet ./...

all: fmt vet build