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

2. **OmniFocus Client** (`internal/omnifocus/`)
   - `client.go`: Go wrapper that executes JXA scripts via `osascript`
   - `types.go`: Data structures for OmniFocus entities (Project, Task, Tag)
   - Handles JSON marshaling/unmarshaling between Go and JXA

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
  main.go                   # Server setup and tool registration
internal/omnifocus/         # OmniFocus client library
  client.go                 # Executes JXA scripts, handles I/O
  types.go                  # Go types for OmniFocus entities
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

### Testing JXA Scripts
Scripts can be tested independently:
```bash
osascript scripts/list_projects.jxa
osascript scripts/list_tasks.jxa
osascript scripts/create_task.jxa '{"name":"Test Task"}'
```

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

## Important Implementation Notes

### Script Path Resolution
The client finds JXA scripts relative to the Go source file location using `runtime.Caller()`. This works for both `go run` and compiled binaries as long as the scripts directory is in the expected location relative to the binary.

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
