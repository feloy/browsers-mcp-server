package mcp

import (
	"bytes"
	"context"
	"slices"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"k8s.io/klog/v2"

	"github.com/feloy/browsers-mcp-server/pkg/config"
	"github.com/feloy/browsers-mcp-server/pkg/version"
)

type Configuration struct {
	Profile Profile

	StaticConfig *config.StaticConfig
}

func (c *Configuration) isToolApplicable(tool server.ServerTool) bool {
	if c.StaticConfig.EnabledTools != nil && !slices.Contains(c.StaticConfig.EnabledTools, tool.Tool.Name) {
		return false
	}
	if c.StaticConfig.DisabledTools != nil && slices.Contains(c.StaticConfig.DisabledTools, tool.Tool.Name) {
		return false
	}
	return true
}

type Server struct {
	configuration *Configuration
	server        *server.MCPServer
	enabledTools  []string
}

func NewServer(configuration Configuration) (*Server, error) {
	var serverOptions []server.ServerOption
	serverOptions = append(serverOptions,
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithToolCapabilities(true),
		server.WithLogging(),
		server.WithToolHandlerMiddleware(toolCallLoggingMiddleware),
	)

	s := &Server{
		configuration: &configuration,
		server: server.NewMCPServer(
			version.BinaryName,
			version.Version,
			serverOptions...,
		),
	}

	applicableTools := make([]server.ServerTool, 0)
	for _, tool := range s.configuration.Profile.GetTools(s) {
		if !s.configuration.isToolApplicable(tool) {
			continue
		}
		applicableTools = append(applicableTools, tool)
		s.enabledTools = append(s.enabledTools, tool.Tool.Name)
	}
	s.server.SetTools(applicableTools...)

	return s, nil
}

func (s *Server) ServeStdio() error {
	return server.ServeStdio(s.server)
}

func (s *Server) GetEnabledTools() []string {
	return s.enabledTools
}

func NewTextResult(content string, err error) *mcp.CallToolResult {
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: err.Error(),
				},
			},
		}
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: content,
			},
		},
	}
}

func NewStructuredResult(content string, structuredContent any, err error) *mcp.CallToolResult {
	result := NewTextResult(content, err)
	result.StructuredContent = structuredContent
	return result
}

func toolCallLoggingMiddleware(next server.ToolHandlerFunc) server.ToolHandlerFunc {
	return func(ctx context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		klog.V(5).Infof("mcp tool call: %s(%v)", ctr.Params.Name, ctr.Params.Arguments)
		if ctr.Header != nil {
			buffer := bytes.NewBuffer(make([]byte, 0))
			if err := ctr.Header.WriteSubset(buffer, map[string]bool{"Authorization": true, "authorization": true}); err == nil {
				klog.V(7).Infof("mcp tool call headers: %s", buffer)
			}
		}
		return next(ctx, ctr)
	}
}
