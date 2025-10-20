.PHONY: build clean install test run validate-jxa

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

# Validate JXA script syntax
validate-jxa:
	@echo "Validating JXA script syntax..."
	@FAILED=0; \
	for script in scripts/*.jxa; do \
		echo "Checking syntax: $$script"; \
		if osacompile -l JavaScript -o /tmp/test.scpt "$$script" 2>&1; then \
			echo "✓ $$script: syntax OK"; \
			rm -f /tmp/test.scpt; \
		else \
			echo "✗ $$script: syntax error detected"; \
			FAILED=1; \
		fi; \
	done; \
	if [ $$FAILED -eq 1 ]; then \
		echo ""; \
		echo "JXA syntax validation failed"; \
		exit 1; \
	fi; \
	echo ""; \
	echo "All JXA scripts passed syntax validation"

# Run all tests
test: validate-jxa
	go test ./...
