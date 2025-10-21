package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/conall/mcp-omnifocus/internal/omnifocus"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const (
	serverName    = "mcp-omnifocus"
	serverVersion = "0.1.0"
)

func main() {
	// Define command line flags
	scriptsPath := flag.String("scripts", "", "Path to the JXA scripts directory (if not specified, auto-detection is used)")
	cacheTTL := flag.Int("cache-ttl", 30, "Cache TTL in seconds (0 to disable caching)")
	flag.Parse()

	// Check for environment variable override
	cacheTTLSeconds := *cacheTTL
	if envTTL := os.Getenv("MCP_OMNIFOCUS_CACHE_TTL"); envTTL != "" {
		if ttl, err := strconv.Atoi(envTTL); err == nil {
			cacheTTLSeconds = ttl
		}
	}

	// Create OmniFocus client with caching
	ttlDuration := time.Duration(cacheTTLSeconds) * time.Second

	// Create a temporary client to get the scripts directory
	tempClient := omnifocus.NewClient()
	scriptsDir := *scriptsPath
	if scriptsDir == "" {
		scriptsDir = tempClient.GetScriptsDir()
	}

	ofClient := omnifocus.NewClientWithCache(scriptsDir, ttlDuration)

	// Log cache configuration
	if cacheTTLSeconds > 0 {
		log.Printf("Cache enabled with TTL: %d seconds", cacheTTLSeconds)
	} else {
		log.Printf("Cache disabled")
	}

	// Create MCP server
	s := server.NewMCPServer(
		serverName,
		serverVersion,
	)

	// Register tools
	registerTools(s, ofClient)

	// Start server with stdio transport
	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func registerTools(s *server.MCPServer, ofClient *omnifocus.Client) {
	// List Projects Tool
	listProjectsTool := mcp.NewTool("list_projects",
		mcp.WithDescription("List all projects in OmniFocus"),
		mcp.WithString("filter",
			mcp.Description("Optional filter for project status (active, on-hold, completed, dropped)"),
		),
	)

	s.AddTool(listProjectsTool, func(args map[string]interface{}) (*mcp.CallToolResult, error) {
		projects, err := ofClient.ListProjects()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list projects: %v", err)), nil
		}

		// Apply filter if provided
		if filterVal, ok := args["filter"].(string); ok && filterVal != "" {
			filtered := []omnifocus.Project{}
			for _, p := range projects {
				if p.Status == filterVal {
					filtered = append(filtered, p)
				}
			}
			projects = filtered
		}

		result, _ := json.MarshalIndent(projects, "", "  ")
		return mcp.NewToolResultText(string(result)), nil
	})

	// List Tasks Tool
	listTasksTool := mcp.NewTool("list_tasks",
		mcp.WithDescription("List tasks in OmniFocus, optionally filtered by project"),
		mcp.WithString("project_id",
			mcp.Description("Optional project ID to filter tasks"),
		),
	)

	s.AddTool(listTasksTool, func(args map[string]interface{}) (*mcp.CallToolResult, error) {
		projectID := ""
		if pid, ok := args["project_id"].(string); ok {
			projectID = pid
		}

		tasks, err := ofClient.ListTasks(projectID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list tasks: %v", err)), nil
		}

		result, _ := json.MarshalIndent(tasks, "", "  ")
		return mcp.NewToolResultText(string(result)), nil
	})

	// List Tags Tool
	listTagsTool := mcp.NewTool("list_tags",
		mcp.WithDescription("List all tags in OmniFocus"),
	)

	s.AddTool(listTagsTool, func(args map[string]interface{}) (*mcp.CallToolResult, error) {
		tags, err := ofClient.ListTags()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list tags: %v", err)), nil
		}

		result, _ := json.MarshalIndent(tags, "", "  ")
		return mcp.NewToolResultText(string(result)), nil
	})

	// Create Task Tool
	createTaskTool := mcp.NewTool("create_task",
		mcp.WithDescription("Create a new task in OmniFocus"),
		mcp.WithString("name",
			mcp.Description("Task name (required)"),
			mcp.Required(),
		),
		mcp.WithString("note",
			mcp.Description("Task note/description"),
		),
		mcp.WithString("project_id",
			mcp.Description("Project ID to add task to (if not provided, adds to inbox)"),
		),
		mcp.WithString("due_date",
			mcp.Description("Due date in ISO 8601 format (e.g., 2024-12-31T23:59:59Z)"),
		),
		mcp.WithBoolean("flagged",
			mcp.Description("Whether to flag the task"),
		),
		mcp.WithNumber("estimated_minutes",
			mcp.Description("Estimated time in minutes"),
		),
		mcp.WithString("tags",
			mcp.Description("Comma-separated list of tag names"),
		),
	)

	s.AddTool(createTaskTool, func(args map[string]interface{}) (*mcp.CallToolResult, error) {
		req := omnifocus.CreateTaskRequest{
			Name: args["name"].(string),
		}

		if note, ok := args["note"].(string); ok {
			req.Note = note
		}
		if projectID, ok := args["project_id"].(string); ok {
			req.ProjectID = projectID
		}
		if dueDate, ok := args["due_date"].(string); ok {
			req.DueDate = dueDate
		}
		if flagged, ok := args["flagged"].(bool); ok {
			req.Flagged = flagged
		}
		if estimatedMinutes, ok := args["estimated_minutes"].(float64); ok {
			minutes := int(estimatedMinutes)
			req.EstimatedMinutes = minutes
		}
		if tagsStr, ok := args["tags"].(string); ok && tagsStr != "" {
			// Simple split by comma
			tags := []string{}
			current := ""
			for _, char := range tagsStr {
				if char == ',' {
					if current != "" {
						tags = append(tags, current)
						current = ""
					}
				} else if char != ' ' || current != "" {
					current += string(char)
				}
			}
			if current != "" {
				tags = append(tags, current)
			}
			req.Tags = tags
		}

		result, err := ofClient.CreateTask(req)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create task: %v", err)), nil
		}

		resultJSON, _ := json.MarshalIndent(result, "", "  ")
		return mcp.NewToolResultText(string(resultJSON)), nil
	})

	// Create Project Tool
	createProjectTool := mcp.NewTool("create_project",
		mcp.WithDescription("Create a new project in OmniFocus"),
		mcp.WithString("name",
			mcp.Description("Project name (required)"),
			mcp.Required(),
		),
		mcp.WithString("note",
			mcp.Description("Project note/description"),
		),
		mcp.WithString("status",
			mcp.Description("Project status (active, on-hold, completed, dropped)"),
		),
		mcp.WithString("tags",
			mcp.Description("Comma-separated list of tag names"),
		),
	)

	s.AddTool(createProjectTool, func(args map[string]interface{}) (*mcp.CallToolResult, error) {
		req := omnifocus.CreateProjectRequest{
			Name: args["name"].(string),
		}

		if note, ok := args["note"].(string); ok {
			req.Note = note
		}
		if status, ok := args["status"].(string); ok {
			req.Status = status
		}
		if tagsStr, ok := args["tags"].(string); ok && tagsStr != "" {
			tags := []string{}
			current := ""
			for _, char := range tagsStr {
				if char == ',' {
					if current != "" {
						tags = append(tags, current)
						current = ""
					}
				} else if char != ' ' || current != "" {
					current += string(char)
				}
			}
			if current != "" {
				tags = append(tags, current)
			}
			req.Tags = tags
		}

		result, err := ofClient.CreateProject(req)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create project: %v", err)), nil
		}

		resultJSON, _ := json.MarshalIndent(result, "", "  ")
		return mcp.NewToolResultText(string(resultJSON)), nil
	})

	// Update Task Tool
	updateTaskTool := mcp.NewTool("update_task",
		mcp.WithDescription("Update an existing task in OmniFocus"),
		mcp.WithString("id",
			mcp.Description("Task ID (required)"),
			mcp.Required(),
		),
		mcp.WithString("name",
			mcp.Description("New task name"),
		),
		mcp.WithString("note",
			mcp.Description("New task note"),
		),
		mcp.WithBoolean("completed",
			mcp.Description("Mark task as completed or incomplete"),
		),
		mcp.WithBoolean("flagged",
			mcp.Description("Flag or unflag the task"),
		),
		mcp.WithString("due_date",
			mcp.Description("New due date in ISO 8601 format (or null to remove)"),
		),
		mcp.WithNumber("estimated_minutes",
			mcp.Description("New estimated time in minutes"),
		),
	)

	s.AddTool(updateTaskTool, func(args map[string]interface{}) (*mcp.CallToolResult, error) {
		req := omnifocus.UpdateTaskRequest{
			ID: args["id"].(string),
		}

		if name, ok := args["name"].(string); ok {
			req.Name = &name
		}
		if note, ok := args["note"].(string); ok {
			req.Note = &note
		}
		if completed, ok := args["completed"].(bool); ok {
			req.Completed = &completed
		}
		if flagged, ok := args["flagged"].(bool); ok {
			req.Flagged = &flagged
		}
		if dueDate, ok := args["due_date"].(string); ok {
			req.DueDate = &dueDate
		}
		if estimatedMinutes, ok := args["estimated_minutes"].(float64); ok {
			minutes := int(estimatedMinutes)
			req.EstimatedMinutes = &minutes
		}

		result, err := ofClient.UpdateTask(req)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to update task: %v", err)), nil
		}

		resultJSON, _ := json.MarshalIndent(result, "", "  ")
		return mcp.NewToolResultText(string(resultJSON)), nil
	})

	// Complete Task Tool
	completeTaskTool := mcp.NewTool("complete_task",
		mcp.WithDescription("Mark a task as complete in OmniFocus"),
		mcp.WithString("id",
			mcp.Description("Task ID (required)"),
			mcp.Required(),
		),
	)

	s.AddTool(completeTaskTool, func(args map[string]interface{}) (*mcp.CallToolResult, error) {
		taskID := args["id"].(string)

		result, err := ofClient.CompleteTask(taskID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to complete task: %v", err)), nil
		}

		resultJSON, _ := json.MarshalIndent(result, "", "  ")
		return mcp.NewToolResultText(string(resultJSON)), nil
	})
}
