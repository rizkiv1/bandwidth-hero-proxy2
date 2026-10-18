.PHONY: all build test run deps clean

# Name of the binary
BINARY_NAME=bandwidth-hero-proxy

all: test build

deps:
	@echo "Checking for libvips..."
	@pkg-config --exists vips || (echo "ERROR: libvips is not installed. Please install libvips-dev (Ubuntu/Debian) or vips-dev (Alpine) before compiling." && exit 1)
	@echo "libvips found."
	go mod download

build: deps
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) cmd/server/main.go

test: deps
	@echo "Running tests..."
	go test -v ./...

run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BINARY_NAME)

clean:
	@echo "Cleaning up..."
	go clean
	rm -f $(BINARY_NAME)
