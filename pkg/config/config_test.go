package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadConfigMissingFile(t *testing.T) {
	config, err := ReadConfig("non-existent-config.toml")
	t.Run("returns error for missing file", func(t *testing.T) {
		if err == nil {
			t.Fatal("Expected error for missing file, got nil")
		}
		if config != nil {
			t.Fatalf("Expected nil config for missing file, got %v", config)
		}
	})
}

func TestReadConfigInvalid(t *testing.T) {
	invalidConfigPath := writeConfig(t, `log_level = 1
port = "9999`)

	config, err := ReadConfig(invalidConfigPath)
	t.Run("returns error for invalid file", func(t *testing.T) {
		if err == nil {
			t.Fatal("Expected error for invalid file, got nil")
		}
		if config != nil {
			t.Fatalf("Expected nil config for invalid file, got %v", config)
		}
	})
	t.Run("error message contains toml error with line number", func(t *testing.T) {
		expectedError := "toml: line 2"
		if err != nil && !strings.HasPrefix(err.Error(), expectedError) {
			t.Fatalf("Expected error message '%s' to contain line number, got %v", expectedError, err)
		}
	})
}

func TestReadConfigValid(t *testing.T) {
	validConfigPath := writeConfig(t, `
log_level = 1
port = "9999"
sse_base_url = "https://example.com"

enabled_tools = ["tool1", "tool2"]
disabled_tools = ["tool3", "tool4"]
`)

	config, err := ReadConfig(validConfigPath)
	t.Run("reads and unmarshalls file", func(t *testing.T) {
		if err != nil {
			t.Fatalf("ReadConfig returned an error for a valid file: %v", err)
		}
		if config == nil {
			t.Fatal("ReadConfig returned a nil config for a valid file")
		}
	})
	t.Run("log_level parsed correctly", func(t *testing.T) {
		if config.LogLevel != 1 {
			t.Fatalf("Unexpected log level: %v", config.LogLevel)
		}
	})
	t.Run("enabled_tools parsed correctly", func(t *testing.T) {
		if len(config.EnabledTools) != 2 {
			t.Fatalf("Unexpected enabled tools: %v", config.EnabledTools)

		}
		for i, tool := range []string{"tool1", "tool2"} {
			if config.EnabledTools[i] != tool {
				t.Errorf("Expected enabled tool %d to be %s, got %s", i, tool, config.EnabledTools[i])
			}
		}
	})
	t.Run("disabled_tools parsed correctly", func(t *testing.T) {
		if len(config.DisabledTools) != 2 {
			t.Fatalf("Unexpected disabled tools: %v", config.DisabledTools)
		}
		for i, tool := range []string{"tool3", "tool4"} {
			if config.DisabledTools[i] != tool {
				t.Errorf("Expected disabled tool %d to be %s, got %s", i, tool, config.DisabledTools[i])
			}
		}
	})
}

func writeConfig(t *testing.T, content string) string {
	t.Helper()
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "config.toml")
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write config file %s: %v", path, err)
	}
	return path
}
