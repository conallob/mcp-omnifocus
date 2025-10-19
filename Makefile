.PHONY: build clean install test run

# Build the MCP server
build:
	go build -o bin/mcp-omnifocus ./cmd/mcp-omnifocus

# Install dependencies
deps:
	go mod download
	go mod tidy

# Clean build artifacts
clean:
	rm -rf bin/

# Install the binary to GOPATH/bin
install:
	go install ./cmd/mcp-omnifocus

# Run the server (for testing)
run:
	go run ./cmd/mcp-omnifocus

# Build for release
release:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/mcp-omnifocus ./cmd/mcp-omnifocus
