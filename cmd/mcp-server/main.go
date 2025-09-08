package main

import (
	"os"

	"github.com/spf13/pflag"

	"github.com/feloy/mcp-server/pkg/genericiooptions"
	"github.com/feloy/mcp-server/pkg/mcp-server/cmd"
)

func main() {
	flags := pflag.NewFlagSet("mcp-server", pflag.ExitOnError)
	pflag.CommandLine = flags

	root := cmd.NewMCPServer(genericiooptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr})
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
