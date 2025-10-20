# mcp-omnifocus

[![Build](https://github.com/conall/mcp-omnifocus/actions/workflows/build.yml/badge.svg)](https://github.com/conall/mcp-omnifocus/actions/workflows/build.yml)
[![Release](https://github.com/conall/mcp-omnifocus/actions/workflows/release.yml/badge.svg)](https://github.com/conall/mcp-omnifocus/actions/workflows/release.yml)
[![GitHub release](https://img.shields.io/github/v/release/conall/mcp-omnifocus)](https://github.com/conall/mcp-omnifocus/releases/latest)

An MCP (Model Context Protocol) server for interacting with OmniFocus Pro, enabling AI assistants to read and manage your OmniFocus projects and tasks through the OmniFocus automation API.

## Features

- **Read Operations**
  - List all projects with their status and metadata
  - List tasks (all tasks or filtered by project)
  - List all tags

- **Write Operations**
  - Create new tasks (in inbox or specific projects)
  - Create new projects
  - Update existing tasks (name, note, status, due date, etc.)
  - Complete tasks
  - Add tags to tasks and projects

## Requirements

- macOS (OmniFocus is macOS/iOS only)
- OmniFocus Pro installed and accessible
- Go 1.21 or later
- System permissions for automation (macOS will prompt on first use)

## Installation

### Option 1: Download Pre-built Binary (Recommended)

1. Download the latest release for your platform from the [Releases page](https://github.com/conall/mcp-omnifocus/releases)
2. Extract the archive:
   ```bash
   tar -xzf mcp-omnifocus_*_darwin_*.tar.gz
   ```
3. Move the binary to your desired location:
   ```bash
   mkdir -p ~/mcp-servers/omnifocus
   mv mcp-omnifocus ~/mcp-servers/omnifocus/
   mv scripts ~/mcp-servers/omnifocus/
   ```

### Option 2: Build from Source

1. Clone this repository:
   ```bash
   git clone https://github.com/conall/mcp-omnifocus.git
   cd mcp-omnifocus
   ```

2. Install dependencies:
   ```bash
   make deps
   ```

3. Build the server:
   ```bash
   make build
   ```

   This creates the binary at `bin/mcp-omnifocus`

## Configuration

To use this MCP server with Claude Desktop or other MCP clients, add it to your configuration file:

### Claude Desktop

Edit your Claude Desktop config file (usually at `~/Library/Application Support/Claude/claude_desktop_config.json`):

**If you downloaded the pre-built binary:**
```json
{
  "mcpServers": {
    "omnifocus": {
      "command": "/Users/YOUR_USERNAME/mcp-servers/omnifocus/mcp-omnifocus",
      "args": ["-scripts", "/Users/YOUR_USERNAME/mcp-servers/omnifocus/scripts"]
    }
  }
}
```

**If you built from source:**
```json
{
  "mcpServers": {
    "omnifocus": {
      "command": "/absolute/path/to/mcp-omnifocus/bin/mcp-omnifocus",
      "args": ["-scripts", "/absolute/path/to/mcp-omnifocus/scripts"]
    }
  }
}
```

Replace the paths with the actual locations where you installed or built the binary and scripts.

**Note:** The `-scripts` argument is optional. If not provided, the server will attempt to auto-detect the scripts directory. However, explicitly specifying the path is recommended for reliability.

## Available Tools

### Read Tools

- **list_projects**: List all projects in OmniFocus
  - Optional `filter` parameter for project status (active, on-hold, completed, dropped)

- **list_tasks**: List tasks in OmniFocus
  - Optional `project_id` parameter to filter tasks by project

- **list_tags**: List all tags in OmniFocus

### Write Tools

- **create_task**: Create a new task
  - Required: `name`
  - Optional: `note`, `project_id`, `due_date`, `flagged`, `estimated_minutes`, `tags`

- **create_project**: Create a new project
  - Required: `name`
  - Optional: `note`, `status`, `tags`

- **update_task**: Update an existing task
  - Required: `id`
  - Optional: `name`, `note`, `completed`, `flagged`, `due_date`, `estimated_minutes`

- **complete_task**: Mark a task as complete
  - Required: `id`

## Architecture

The server is built in Go and uses:
- **JXA (JavaScript for Automation)** to interact with OmniFocus's automation API
- **mcp-go** SDK for implementing the MCP protocol
- **stdio transport** for communication with MCP clients

### Project Structure

```
.
├── cmd/mcp-omnifocus/     # Main server executable
├── internal/omnifocus/     # OmniFocus client library
│   ├── client.go          # Go wrapper for JXA scripts
│   └── types.go           # Data structures
├── scripts/               # JXA scripts for OmniFocus automation
│   ├── list_projects.jxa
│   ├── list_tasks.jxa
│   ├── create_task.jxa
│   └── ...
└── Makefile              # Build configuration
```

## Development

### Building
```bash
make build
```

### Running for testing
```bash
make run
```

### Clean build artifacts
```bash
make clean
```

## macOS Permissions

On first use, macOS will prompt you to grant automation permissions to:
- `osascript` (to run the JXA scripts)
- Access to OmniFocus

These permissions are required for the server to function.

## License

See LICENSE file for details.
