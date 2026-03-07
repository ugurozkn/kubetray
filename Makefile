.PHONY: build test clean install help

BINARY_NAME=kubetray
GOCMD=go

## build: Build the binary
build:
	$(GOCMD) build -o $(BINARY_NAME) .

## test: Run unit tests
test:
	$(GOCMD) test -v -race ./...

## clean: Clean build artifacts
clean:
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html

## install: Install binary to /usr/local/bin
install: build
	sudo mv $(BINARY_NAME) /usr/local/bin/

## deps: Download dependencies
deps:
	$(GOCMD) mod download
	$(GOCMD) mod tidy

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'
