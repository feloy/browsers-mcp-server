package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"k8s.io/klog/v2"

	"github.com/feloy/browsers-mcp-server/pkg/config"
	"github.com/feloy/browsers-mcp-server/pkg/genericiooptions"
	"github.com/feloy/browsers-mcp-server/pkg/mcp"
	"github.com/feloy/browsers-mcp-server/pkg/version"
)

var (
	long     = "Model Context Protocol (MCP) server"
	examples = `
# show this help
mcp-server -h

# shows version information
mcp-server --version

# start STDIO server
mcp-server`
)

type MCPServerOptions struct {
	Version  bool
	LogLevel int

	Profile string

	ConfigPath   string
	StaticConfig *config.StaticConfig

	genericiooptions.IOStreams
	Logger
}

func NewMCPServerOptions(streams genericiooptions.IOStreams) *MCPServerOptions {
	return &MCPServerOptions{
		IOStreams:    streams,
		Profile:      "full",
		StaticConfig: &config.StaticConfig{},
	}
}

func NewMCPServer(streams genericiooptions.IOStreams) *cobra.Command {
	o := NewMCPServerOptions(streams)
	cmd := &cobra.Command{
		Use:     "mcp-server [command] [options]",
		Short:   "Model Context Protocol (MCP) server",
		Long:    long,
		Example: examples,
		RunE: func(c *cobra.Command, args []string) error {
			o.initLogger()
			defer o.disposeLogger()

			if err := o.Complete(c); err != nil {
				return err
			}
			if err := o.Validate(); err != nil {
				return err
			}
			if err := o.Run(); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&o.Version, "version", o.Version, "Print version information and quit")
	cmd.Flags().StringVar(&o.ConfigPath, "config", o.ConfigPath, "Path of the config file. Each profile has its set of defaults.")
	o.initLoggerFlags(cmd)
	return cmd
}

func (m *MCPServerOptions) Complete(cmd *cobra.Command) error {
	if m.ConfigPath != "" {
		cnf, err := config.ReadConfig(m.ConfigPath)
		if err != nil {
			return err
		}
		m.StaticConfig = cnf
	}

	return nil
}

func (m *MCPServerOptions) Validate() error {
	return nil
}

func (m *MCPServerOptions) Run() error {
	profile := mcp.ProfileFromString(m.Profile)

	klog.V(1).Info("Starting mcp-server")
	klog.V(1).Infof(" - Config: %s", m.ConfigPath)

	if m.Version {
		_, _ = fmt.Fprintf(m.Out, "%s\n", version.Version)
		return nil
	}

	mcpServer, err := mcp.NewServer(mcp.Configuration{
		Profile:      profile,
		StaticConfig: m.StaticConfig,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize MCP server: %w", err)
	}

	if err := mcpServer.ServeStdio(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	return nil
}
