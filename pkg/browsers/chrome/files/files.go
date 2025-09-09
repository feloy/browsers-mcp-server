package files

import (
	"os"
	"path/filepath"

	"github.com/feloy/browsers-mcp-server/pkg/system"
)

func getUserDataDirecory() string {
	if system.Os == "darwin" {
		return filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "Google", "Chrome")
	}
	if system.Os == "windows" {
		return filepath.Join(os.Getenv("APPDATA"), "Google", "Chrome", "User Data") // TODO check APPDATA on Windows platform
	}
	if system.Os == "linux" {
		return filepath.Join(os.Getenv("HOME"), ".config", "google-chrome")
	}
	return ""
}
