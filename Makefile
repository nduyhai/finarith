.PHONY: deps test lint clean

# Default target
all: deps test

# Download dependencies
deps:
	go mod download
	go mod tidy

# Run tests
test:
	go test -v ./...

# Run linting
lint:
	go vet ./...
	go fmt ./...

# Clean build artifacts
clean:
	go clean
	rm -f coverage.out

# Show help
help:
	@echo "Available targets:"
	@echo "  deps   - Download dependencies"
	@echo "  test   - Run tests"
	@echo "  lint   - Run linting tools"
	@echo "  clean  - Clean build artifacts"
	@echo "  all    - Run deps and test (default)"
	@echo "  help   - Show this help message"