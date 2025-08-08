package tools

import (
	"context"
	"encoding/json"
	"os"

	"github.com/kubeflow/model-registry/mcp/internal/defaults"
	"github.com/kubeflow/model-registry/mcp/internal/mr"
	"github.com/kubeflow/model-registry/pkg/openapi"
	mcpgolang "github.com/metoro-io/mcp-golang"
)

type RegisterModelHandlerArguments struct {
	ModelName string `json:"model_name" jsonschema:"required,description=The name of the model to register"`
}

func RegisterModelHandler(ctx context.Context, arguments RegisterModelHandlerArguments) (*mcpgolang.ToolResponse, error) {
	client := mr.GetInstance().GetClient(os.Getenv(defaults.ModelRegistryURL))

	createModel := openapi.RegisteredModelCreate{
		Name: arguments.ModelName,
	}

	model, _, err := client.ModelRegistryServiceAPI.CreateRegisteredModel(ctx).RegisteredModelCreate(createModel).Execute()
	if err != nil {
		return nil, err
	}

	outputJson, err := json.Marshal(model)
	if err != nil {
		return nil, err
	}

	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(string(outputJson))), nil
}
