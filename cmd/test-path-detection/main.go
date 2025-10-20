package main

import (
	"fmt"
	"os"

	"github.com/conall/mcp-omnifocus/internal/omnifocus"
)

func main() {
	fmt.Println("Testing OmniFocus client script path detection...")
	fmt.Println()

	// Create a new client
	client := omnifocus.NewClient()

	fmt.Printf("✓ Client created successfully\n")
	fmt.Printf("  Scripts directory detected at: %s\n", client.GetScriptsDir())

	// Verify the scripts directory exists
	scriptsDir := client.GetScriptsDir()
	if _, err := os.Stat(scriptsDir); os.IsNotExist(err) {
		fmt.Printf("✗ ERROR: Scripts directory does not exist: %s\n", scriptsDir)
		os.Exit(1)
	}

	// Check for required scripts
	requiredScripts := []string{
		"list_projects.jxa",
		"list_tasks.jxa",
		"list_tags.jxa",
		"create_task.jxa",
		"create_project.jxa",
		"update_task.jxa",
		"complete_task.jxa",
	}

	fmt.Println()
	fmt.Println("Checking for required scripts:")
	allFound := true
	for _, script := range requiredScripts {
		scriptPath := fmt.Sprintf("%s/%s", scriptsDir, script)
		if _, err := os.Stat(scriptPath); err == nil {
			fmt.Printf("  ✓ %s\n", script)
		} else {
			fmt.Printf("  ✗ %s (not found)\n", script)
			allFound = false
		}
	}

	fmt.Println()
	if allFound {
		fmt.Println("✓ All script path detection tests passed!")
		os.Exit(0)
	} else {
		fmt.Println("✗ Some scripts were not found")
		os.Exit(1)
	}
}
