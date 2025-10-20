package omnifocus

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// Client provides methods to interact with OmniFocus
type Client struct {
	scriptsDir string
}

// NewClient creates a new OmniFocus client
func NewClient() *Client {
	scriptsDir := findScriptsDir()
	return &Client{
		scriptsDir: scriptsDir,
	}
}

// GetScriptsDir returns the path to the scripts directory
// This is primarily used for testing and debugging
func (c *Client) GetScriptsDir() string {
	return c.scriptsDir
}

// findScriptsDir attempts to locate the scripts directory in multiple locations
func findScriptsDir() string {
	// Get the executable path
	execPath, err := os.Executable()
	if err != nil {
		execPath = ""
	} else {
		// Resolve symlinks (important for Homebrew installations)
		execPath, err = filepath.EvalSymlinks(execPath)
		if err != nil {
			execPath = ""
		}
	}

	// List of candidate paths to check, in order of preference
	var candidates []string

	if execPath != "" {
		execDir := filepath.Dir(execPath)

		// 1. Scripts directory next to the binary (release package layout)
		candidates = append(candidates, filepath.Join(execDir, "scripts"))

		// 2. Scripts directory one level up (for bin/mcp-omnifocus structure)
		candidates = append(candidates, filepath.Join(execDir, "..", "scripts"))

		// 3. Homebrew installation path (share/mcp-omnifocus/)
		candidates = append(candidates, filepath.Join(execDir, "..", "share", "mcp-omnifocus", "scripts"))
	}

	// 4. Relative to the Go source file (development mode with go run)
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		projectRoot := filepath.Join(filepath.Dir(filename), "..", "..")
		candidates = append(candidates, filepath.Join(projectRoot, "scripts"))
	}

	// 5. Check current working directory (fallback)
	if cwd, err := os.Getwd(); err == nil {
		candidates = append(candidates, filepath.Join(cwd, "scripts"))
	}

	// Try each candidate path
	for _, dir := range candidates {
		absDir, err := filepath.Abs(dir)
		if err != nil {
			continue
		}

		// Check if the directory exists and contains at least one .jxa file
		if isValidScriptsDir(absDir) {
			return absDir
		}
	}

	// If nothing found, fall back to the development layout
	_, filename, _, _ = runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(filename), "..", "..")
	return filepath.Join(projectRoot, "scripts")
}

// isValidScriptsDir checks if a directory exists and contains .jxa files
func isValidScriptsDir(dir string) bool {
	// Check if directory exists
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		return false
	}

	// Check for at least one required script file
	requiredScripts := []string{"list_projects.jxa", "list_tasks.jxa", "create_task.jxa"}
	for _, script := range requiredScripts {
		scriptPath := filepath.Join(dir, script)
		if _, err := os.Stat(scriptPath); err == nil {
			return true
		}
	}

	return false
}

// executeJXA executes a JXA script and returns the output
func (c *Client) executeJXA(scriptName string, args ...string) ([]byte, error) {
	scriptPath := filepath.Join(c.scriptsDir, scriptName)

	cmdArgs := append([]string{scriptPath}, args...)
	cmd := exec.Command("osascript", cmdArgs...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to execute %s: %w - %s", scriptName, err, string(output))
	}

	return output, nil
}

// ListProjects retrieves all projects from OmniFocus
func (c *Client) ListProjects() ([]Project, error) {
	output, err := c.executeJXA("list_projects.jxa")
	if err != nil {
		return nil, err
	}

	var projects []Project
	if err := json.Unmarshal(output, &projects); err != nil {
		return nil, fmt.Errorf("failed to parse projects: %w", err)
	}

	return projects, nil
}

// ListTasks retrieves tasks from OmniFocus, optionally filtered by project ID
func (c *Client) ListTasks(projectID string) ([]Task, error) {
	var output []byte
	var err error

	if projectID != "" {
		output, err = c.executeJXA("list_tasks.jxa", projectID)
	} else {
		output, err = c.executeJXA("list_tasks.jxa")
	}

	if err != nil {
		return nil, err
	}

	var tasks []Task
	if err := json.Unmarshal(output, &tasks); err != nil {
		return nil, fmt.Errorf("failed to parse tasks: %w", err)
	}

	return tasks, nil
}

// ListTags retrieves all tags from OmniFocus
func (c *Client) ListTags() ([]Tag, error) {
	output, err := c.executeJXA("list_tags.jxa")
	if err != nil {
		return nil, err
	}

	var tags []Tag
	if err := json.Unmarshal(output, &tags); err != nil {
		return nil, fmt.Errorf("failed to parse tags: %w", err)
	}

	return tags, nil
}

// CreateTask creates a new task in OmniFocus
func (c *Client) CreateTask(req CreateTaskRequest) (*OperationResult, error) {
	reqJSON, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	output, err := c.executeJXA("create_task.jxa", string(reqJSON))
	if err != nil {
		return nil, err
	}

	var result OperationResult
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}

	if result.Error != "" {
		return &result, fmt.Errorf("OmniFocus error: %s", result.Error)
	}

	return &result, nil
}

// CreateProject creates a new project in OmniFocus
func (c *Client) CreateProject(req CreateProjectRequest) (*OperationResult, error) {
	reqJSON, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	output, err := c.executeJXA("create_project.jxa", string(reqJSON))
	if err != nil {
		return nil, err
	}

	var result OperationResult
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}

	if result.Error != "" {
		return &result, fmt.Errorf("OmniFocus error: %s", result.Error)
	}

	return &result, nil
}

// UpdateTask updates an existing task in OmniFocus
func (c *Client) UpdateTask(req UpdateTaskRequest) (*OperationResult, error) {
	reqJSON, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	output, err := c.executeJXA("update_task.jxa", string(reqJSON))
	if err != nil {
		return nil, err
	}

	var result OperationResult
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}

	if result.Error != "" {
		return &result, fmt.Errorf("OmniFocus error: %s", result.Error)
	}

	return &result, nil
}

// CompleteTask marks a task as complete in OmniFocus
func (c *Client) CompleteTask(taskID string) (*OperationResult, error) {
	output, err := c.executeJXA("complete_task.jxa", taskID)
	if err != nil {
		return nil, err
	}

	var result OperationResult
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}

	if result.Error != "" {
		return &result, fmt.Errorf("OmniFocus error: %s", result.Error)
	}

	return &result, nil
}
