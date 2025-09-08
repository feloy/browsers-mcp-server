package files

import (
	"os"
	"path/filepath"

	"github.com/feloy/mcp-server/pkg/system"
)

func getUserDataDirecory() string {
	if system.Os == "darwin" {
		return filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "Firefox")
	}
	if system.Os == "windows" {
		return filepath.Join(os.Getenv("APPDATA"), "Mozilla", "Firefox") // TODO check APPDATA on Windows platform
	}
	if system.Os == "linux" {
		return filepath.Join(os.Getenv("HOME"), ".mozilla", "firefox")
	}
	return ""
}
