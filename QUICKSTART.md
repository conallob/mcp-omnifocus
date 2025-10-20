# Quick Start Guide

This guide will help you get the MCP OmniFocus server running with Claude Desktop.

## Prerequisites

1. **macOS** with OmniFocus Pro installed
2. **Go 1.21+** installed (check with `go version`)
3. **Claude Desktop** installed

## Installation

### 1. Build the Server

```bash
# Clone or navigate to the repository
cd mcp-omnifocus

# Download dependencies
make deps

# Build the server
make build
```

This creates the executable at `bin/mcp-omnifocus`.

### 2. Configure Claude Desktop

1. Open your Claude Desktop configuration file:
   ```bash
   open ~/Library/Application\ Support/Claude/claude_desktop_config.json
   ```

2. Add the OmniFocus MCP server to your configuration:
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

   **Important**: Replace `/absolute/path/to/mcp-omnifocus` with the actual absolute path to your repository. For example:
   ```json
   {
     "mcpServers": {
       "omnifocus": {
         "command": "/Users/conall/Documents/GitHub/mcp-omnifocus/bin/mcp-omnifocus",
         "args": ["-scripts", "/Users/conall/Documents/GitHub/mcp-omnifocus/scripts"]
       }
     }
   }
   ```

   **Note**: The `-scripts` argument is optional. If omitted, the server will auto-detect the scripts directory, but specifying it explicitly is more reliable.

3. Save the file and restart Claude Desktop.

### 3. Grant Permissions

On first use, macOS will prompt you to grant automation permissions:

1. Open **System Preferences** → **Security & Privacy** → **Privacy** → **Automation**
2. Grant `osascript` permission to control OmniFocus
3. You may also need to grant Terminal or Claude Desktop permission

### 4. Test the Server

Open Claude Desktop and try these example commands:

**List all projects:**
```
Use the list_projects tool to show me all my OmniFocus projects
```

**Create a task:**
```
Create a new task in OmniFocus called "Test MCP integration" with the note "This is a test from Claude"
```

**List tasks:**
```
Show me all my tasks in OmniFocus
```

## Available Tools

Once configured, Claude can use these tools to interact with OmniFocus:

### Read Tools
- `list_projects` - Get all projects (optionally filter by status)
- `list_tasks` - Get all tasks or tasks in a specific project
- `list_tags` - Get all tags

### Write Tools
- `create_task` - Create a new task in inbox or a specific project
- `create_project` - Create a new project
- `update_task` - Update task properties (name, note, status, due date, etc.)
- `complete_task` - Mark a task as complete

## Example Workflows

### 1. Daily Planning
```
Look at all my active projects and incomplete tasks. Help me identify what I should focus on today.
```

### 2. Create Tasks from Notes
```
I need to:
- Review the Q4 budget
- Call Sarah about the new proposal
- Finish the presentation slides

Create these as tasks in OmniFocus with a due date of Friday.
```

### 3. Project Management
```
Create a new project called "Website Redesign" and add these tasks to it:
1. Research design trends
2. Create wireframes
3. Get stakeholder feedback
4. Build prototype
```

## Troubleshooting

### "osascript" is not allowed to control OmniFocus
- Go to System Preferences → Security & Privacy → Automation
- Check the box next to OmniFocus under osascript

### Server not found in Claude Desktop
- Verify the path in `claude_desktop_config.json` is absolute (starts with `/`)
- Make sure the binary exists at that path: `ls -l /path/to/bin/mcp-omnifocus`
- Restart Claude Desktop completely

### JXA Script Errors
- Ensure OmniFocus is running
- Check that scripts have execute permissions: `chmod +x scripts/*.jxa`
- Test scripts manually: `osascript scripts/list_projects.jxa`

### Changes not appearing
- Restart Claude Desktop after modifying the config
- Check Claude Desktop logs for errors
- Rebuild the binary if you modified the Go code: `make build`

## Next Steps

- Read the full [README.md](README.md) for detailed documentation
- Check [CLAUDE.md](CLAUDE.md) for development guidance
- Explore the JXA scripts in `scripts/` directory to understand the automation

## Support

For issues or questions:
- Check the [OmniFocus Automation documentation](https://omni-automation.com/omnifocus/)
- Review the MCP protocol documentation
- Open an issue on GitHub
