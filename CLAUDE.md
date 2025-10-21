# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is an MCP (Model Context Protocol) server written in Go that enables AI assistants to interact with OmniFocus Pro through its automation API. The server uses JXA (JavaScript for Automation) scripts to communicate with OmniFocus and exposes this functionality through the MCP protocol.

## Architecture

### Components

1. **Go MCP Server** (`cmd/mcp-omnifocus/main.go`)
   - Main entry point for the MCP server
   - Uses the `mcp-go` SDK (`github.com/mark3labs/mcp-go`)
   - Implements stdio transport for MCP communication
   - Registers all MCP tools and their handlers
   - Configures caching layer with TTL settings

2. **OmniFocus Client** (`internal/omnifocus/`)
   - `client.go`: Go wrapper that executes JXA scripts via `osascript`
   - `types.go`: Data structures for OmniFocus entities (Project, Task, Tag)
   - `cache.go`: In-memory caching layer with TTL support
   - Handles JSON marshaling/unmarshaling between Go and JXA
   - Implements cache invalidation on write operations

3. **JXA Automation Scripts** (`scripts/*.jxa`)
   - Bridge layer between Go and OmniFocus
   - Each script handles a specific operation (list, create, update)
   - Returns JSON for easy parsing in Go
   - Scripts are executable and can be tested independently

### Data Flow

```
MCP Client → stdio → Go Server → osascript → JXA Script → OmniFocus API
                                                              ↓
MCP Client ← stdio ← Go Server ← JSON output ← JXA Script ← OmniFocus
```

## Project Structure

```
cmd/mcp-omnifocus/          # Main server binary
  main.go                   # Server setup, tool registration, cache config
internal/omnifocus/         # OmniFocus client library
  client.go                 # Executes JXA scripts, handles I/O, caching
  client_test.go           # Tests for client functionality
  types.go                  # Go types for OmniFocus entities
  cache.go                  # In-memory cache with TTL
  cache_test.go            # Tests for caching functionality
scripts/                    # JXA automation scripts
  list_projects.jxa         # Retrieve all projects
  list_tasks.jxa           # Retrieve tasks (all or by project)
  list_tags.jxa            # Retrieve all tags
  create_task.jxa          # Create a new task
  create_project.jxa       # Create a new project
  update_task.jxa          # Update task properties
  complete_task.jxa        # Mark task complete
```

## Common Development Commands

### Building
- `make build` - Build the server binary to `bin/mcp-omnifocus`
- `make deps` - Download and tidy Go dependencies
- `make clean` - Remove build artifacts
- `make release` - Build optimized binary for release

### Running
- `make run` - Run the server directly with `go run`
- `./bin/mcp-omnifocus` - Run the compiled binary
- `./bin/mcp-omnifocus -scripts /path/to/scripts` - Run with explicit scripts path

### Testing
- `make test` - Run all tests (includes Go tests and JXA validation)
- `make validate-jxa` - Validate JXA script syntax only
- `go test ./... -v` - Run Go tests only

### Testing JXA Scripts
Scripts can be tested independently:
```bash
osascript -l JavaScript scripts/list_projects.jxa
osascript -l JavaScript scripts/list_tasks.jxa
osascript -l JavaScript scripts/create_task.jxa '{"name":"Test Task"}'
```

**Important:** Always use the `-l JavaScript` flag when testing JXA scripts with `osascript` to ensure proper syntax parsing.

## MCP Tools

All tools are registered in `cmd/mcp-omnifocus/main.go` in the `registerTools()` function.

### Read Tools
- `list_projects` - Returns all projects (optional status filter)
- `list_tasks` - Returns tasks (optional project_id filter)
- `list_tags` - Returns all tags

### Write Tools
- `create_task` - Creates task (inbox or in project)
- `create_project` - Creates new project
- `update_task` - Updates task properties
- `complete_task` - Marks task complete

## Caching Implementation

The server includes a built-in caching layer to improve performance:

### Cache Architecture

- **Location**: `internal/omnifocus/cache.go`
- **Type**: In-memory cache with TTL (Time To Live)
- **Thread-safety**: Uses `sync.RWMutex` for concurrent access
- **Cleanup**: Automatic background cleanup of expired entries

### Cache Behavior

1. **Read Operations** (cached):
   - `ListProjects()` - Cache key: `"projects:all"`
   - `ListTasks("")` - Cache key: `"tasks:all"`
   - `ListTasks(projectID)` - Cache key: `"tasks:project:<projectID>"`
   - `ListTags()` - Cache key: `"tags:all"`

2. **Write Operations** (invalidate cache):
   - `CreateTask()` - Invalidates all task caches, and project caches if task added to project
   - `CreateProject()` - Invalidates all project caches
   - `UpdateTask()` - Invalidates all task and project caches
   - `CompleteTask()` - Invalidates all task and project caches

3. **Configuration**:
   - Default TTL: 30 seconds
   - Command-line flag: `-cache-ttl <seconds>`
   - Environment variable: `MCP_OMNIFOCUS_CACHE_TTL`
   - Disable caching: Set TTL to 0

### Cache Methods

- `Get(key)` - Retrieve value if exists and not expired
- `Set(key, value)` - Store value with TTL
- `Invalidate(key)` - Remove specific key
- `InvalidateAll()` - Clear entire cache
- `InvalidatePattern(prefix)` - Remove all keys starting with prefix
- `Cleanup()` - Remove expired entries
- `StartCleanupTimer(interval)` - Background cleanup goroutine

## Important Implementation Notes

### Script Path Resolution

The server supports two ways to specify the scripts directory:

#### 1. Command Line Flag (Recommended)
Use the `-scripts` flag to explicitly specify the path:
```bash
./bin/mcp-omnifocus -scripts /path/to/scripts
```

This is the recommended approach for MCP configurations:
```json
{
  "mcpServers": {
    "omnifocus": {
      "command": "/path/to/mcp-omnifocus",
      "args": ["-scripts", "/path/to/scripts"]
    }
  }
}
```

#### 2. Auto-Detection (Fallback)
If the `-scripts` flag is not provided, the client automatically detects the JXA scripts directory by checking multiple locations in this order:
1. Scripts directory next to the binary (release package layout)
2. Scripts directory one level up (for bin/mcp-omnifocus structure)
3. Homebrew installation path (`../share/mcp-omnifocus/scripts/`)
4. Parent directory of the binary (for nested archive extractions)
5. Walking up the directory tree (up to 3 levels)
6. Relative to the Go source file (development mode with `go run`)
7. Current working directory and its parent

This robust path detection works for:
- Development mode (`go run`)
- Compiled binaries in the project
- Extracted release archives
- Homebrew installations
- CI/CD environments

#### Debugging Path Detection
If the scripts are not being found, enable debug logging:
```bash
export MCP_OMNIFOCUS_DEBUG=1
./bin/mcp-omnifocus
```

This will output detailed information about:
- The executable path and its resolution
- All candidate paths being checked
- Which path was ultimately selected

You can also use the test utility:
```bash
go run ./cmd/test-path-detection
```

### JSON Communication
- Go marshals request data to JSON strings
- JXA scripts receive JSON as command-line arguments
- JXA returns JSON strings to stdout
- Go unmarshals responses into typed structs

### Error Handling
- JXA scripts return `{error: "message"}` for errors
- Go client checks for error field in responses
- MCP tools return `mcp.NewToolResultError()` for errors
- All errors include context about which operation failed

### macOS Permissions
First run requires granting automation permissions:
- System Preferences → Security & Privacy → Automation
- Grant `osascript` permission to control OmniFocus

## Adding New Tools

To add a new MCP tool:

1. Create JXA script in `scripts/` that returns JSON
2. Add Go method to `internal/omnifocus/client.go`
3. Add types to `internal/omnifocus/types.go` if needed
4. Register tool in `cmd/mcp-omnifocus/main.go` using `mcp.NewTool()`
5. Implement handler function that calls the client method

Example:
```go
myTool := mcp.NewTool("my_tool",
    mcp.WithDescription("Description of what it does"),
    mcp.WithString("param", mcp.Required(true)),
)
s.AddTool(myTool, func(args map[string]interface{}) (*mcp.CallToolResult, error) {
    result, err := ofClient.MyMethod(args["param"].(string))
    // ... handle error and return result
})
```

## Testing

To test the MCP server locally:
1. Build with `make build`
2. Configure in Claude Desktop or use MCP inspector
3. Check that JXA scripts have execute permissions (`chmod +x scripts/*.jxa`)
4. Ensure OmniFocus is running and has projects/tasks to query

## Continuous Integration

The project uses GitHub Actions for CI/CD. The build workflow (`.github/workflows/build.yml`) runs on every push and pull request to `main`:

1. **Build** - Compiles the Go binary
2. **Validate JXA Syntax** - Uses `osacompile` to check all JXA scripts for syntax errors without executing them
3. **Run Tests** - Executes all Go tests

The JXA syntax validation ensures that script syntax errors are caught early, before they cause runtime failures. This is particularly important because the `-l JavaScript` flag must be used when calling `osascript` for proper syntax parsing.

## Release Process

Releases are automated using GitHub Actions and GoReleaser v2.

### Creating a Release

1. Ensure all changes are committed to `main`
2. Create and push a version tag:
   ```bash
   git tag -a v0.1.0 -m "Release v0.1.0"
   git push origin v0.1.0
   ```
3. GitHub Actions automatically builds and publishes the release

### What Gets Built

- macOS binaries (Intel and Apple Silicon)
- Archives with binary, scripts, and documentation
- Checksums for verification
- GitHub Release with changelog

### Files Included in Release

- `mcp-omnifocus` binary
- `scripts/*.jxa` (JXA automation scripts)
- `README.md`, `LICENSE`, `QUICKSTART.md`
- `config-example.json`

**Important**: The `scripts/` directory must be in the same location as the binary or in a known relative path for the client to find the JXA scripts.

See [RELEASE.md](RELEASE.md) for detailed release instructions.
