package config

import (
	"github.com/BurntSushi/toml"
	"github.com/feloy/browsers-mcp-server/pkg/system"
)

// StaticConfig is the configuration for the server.
// It allows to configure server specific settings and tools to be enabled or disabled.
type StaticConfig struct {
	LogLevel      int      `toml:"log_level,omitempty"`
	EnabledTools  []string `toml:"enabled_tools,omitempty"`
	DisabledTools []string `toml:"disabled_tools,omitempty"`
}

// ReadConfig reads the toml file and returns the StaticConfig.
func ReadConfig(configPath string) (*StaticConfig, error) {
	configData, err := system.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config *StaticConfig
	err = toml.Unmarshal(configData, &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
