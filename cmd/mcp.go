package cmd

import (
	"github.com/kubeflow/model-registry/mcp/cmd"
)

func init() {
	rootCmd.AddCommand(cmd.MCPCmd)
}
