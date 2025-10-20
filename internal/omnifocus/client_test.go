package omnifocus

import (
	"path/filepath"
	"testing"
)

func TestFindScriptsDir(t *testing.T) {
	scriptsDir := findScriptsDir()

	// Check that we got a non-empty path
	if scriptsDir == "" {
		t.Fatal("findScriptsDir returned empty string")
	}

	// Check that the path is absolute
	if !filepath.IsAbs(scriptsDir) {
		t.Errorf("findScriptsDir returned relative path: %s", scriptsDir)
	}

	// Check that it's a valid scripts directory
	if !isValidScriptsDir(scriptsDir) {
		t.Errorf("findScriptsDir returned invalid directory: %s", scriptsDir)
	}

	t.Logf("Found scripts directory at: %s", scriptsDir)
}

func TestIsValidScriptsDir(t *testing.T) {
	tests := []struct {
		name     string
		dir      string
		expected bool
	}{
		{
			name:     "non-existent directory",
			dir:      "/path/that/does/not/exist",
			expected: false,
		},
		{
			name:     "empty directory path",
			dir:      "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidScriptsDir(tt.dir)
			if result != tt.expected {
				t.Errorf("isValidScriptsDir(%q) = %v, want %v", tt.dir, result, tt.expected)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	client := NewClient()

	if client == nil {
		t.Fatal("NewClient returned nil")
	}

	if client.scriptsDir == "" {
		t.Error("Client scriptsDir is empty")
	}

	// Verify the scripts directory is valid
	if !isValidScriptsDir(client.scriptsDir) {
		t.Errorf("Client has invalid scriptsDir: %s", client.scriptsDir)
	}

	t.Logf("Client initialized with scriptsDir: %s", client.scriptsDir)
}

func TestNewClientWithPath(t *testing.T) {
	// Test with a custom path
	customPath := "/custom/path/to/scripts"
	client := NewClientWithPath(customPath)

	if client == nil {
		t.Fatal("NewClientWithPath returned nil")
	}

	if client.scriptsDir != customPath {
		t.Errorf("Client scriptsDir = %s, want %s", client.scriptsDir, customPath)
	}

	t.Logf("Client initialized with custom scriptsDir: %s", client.scriptsDir)

	// Test with the actual scripts directory
	actualScriptsDir := findScriptsDir()
	client2 := NewClientWithPath(actualScriptsDir)

	if client2 == nil {
		t.Fatal("NewClientWithPath returned nil for actual scripts dir")
	}

	if client2.scriptsDir != actualScriptsDir {
		t.Errorf("Client scriptsDir = %s, want %s", client2.scriptsDir, actualScriptsDir)
	}

	// Verify the GetScriptsDir method returns the correct value
	if client2.GetScriptsDir() != actualScriptsDir {
		t.Errorf("GetScriptsDir() = %s, want %s", client2.GetScriptsDir(), actualScriptsDir)
	}

	t.Logf("Client initialized with actual scriptsDir: %s", client2.scriptsDir)
}
