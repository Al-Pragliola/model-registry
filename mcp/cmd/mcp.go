package cmd

import (
	"fmt"
	"os"

	"github.com/kubeflow/model-registry/mcp/internal/defaults"
	"github.com/kubeflow/model-registry/mcp/internal/tools"
	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
	"github.com/spf13/cobra"
)

var catalogCfg = struct {
	ListenAddress    string
	ModelRegistryURL string
}{
	ListenAddress: "0.0.0.0:8080",
}

var MCPCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Model Registry MCP server",
	Long:  `Launch the MCP server for the model registry`,
	RunE:  runMCPServer,
}

func init() {
	MCPCmd.Flags().StringVarP(&catalogCfg.ListenAddress, "listen", "l", catalogCfg.ListenAddress, "Address to listen on")
}

func runMCPServer(cmd *cobra.Command, args []string) error {
	if os.Getenv(defaults.ModelRegistryURL) == "" {
		return fmt.Errorf("%s environment variable is not set", defaults.ModelRegistryURL)
	}

	done := make(chan struct{})

	server := mcp_golang.NewServer(stdio.NewStdioServerTransport())

	for _, tool := range tools.Tools {
		err := server.RegisterTool(tool.Name, tool.Description, tool.Handler)
		if err != nil {
			return fmt.Errorf("error registering tool %s: %v", tool.Name, err)
		}
	}

	err := server.Serve()
	if err != nil {
		panic(err)
	}

	<-done

	return nil
}
