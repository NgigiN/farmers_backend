BINARY_NAME=farmers_backend
MAIN_PATH=./main.go

.PHONY: help build run clean test deps fmt vet docker-build docker-up docker-down docker-restart

help:
	@echo "  make build         - Build the application"
	@echo "  make run           - Run the application"
	@echo "  make test          - Run tests"
	@echo "  make clean         - Remove build artifacts"
	@echo "  make deps          - Download dependencies"
	@echo "  make fmt           - Format code"
	@echo "  make vet           - Run go vet"
	@echo "  make all           - Format, vet, and build"
	@echo "  make docker-build  - Build Docker image"
	@echo "  make docker-up     - Start Docker container"
	@echo "  make docker-down   - Stop Docker container"
	@echo "  make docker-restart - Restart Docker container"

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

docker-build:
	docker-compose build

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-restart: docker-down docker-up