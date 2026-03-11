package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ugurozkn/kubetray/pkg/mcp"
)

var mcpCmd = &cobra.Command{
	Use:    "mcp",
	Short:  "Start MCP server for AI tool integration",
	Long:   `Start a Model Context Protocol (MCP) server over stdin/stdout. This allows AI assistants like Claude to manage your local Kubernetes cluster.`,
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		server := mcp.NewServer(version)
		return server.Run()
	},
}

func init() {
	rootCmd.AddCommand(mcpCmd)
}
