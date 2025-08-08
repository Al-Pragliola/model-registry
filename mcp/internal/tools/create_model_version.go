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

type CreateModelVersionHandlerArguments struct {
	ModelVersionName  string `json:"model_version_name" jsonschema:"required,description=The name of the model version to create"`
	RegisteredModelID string `json:"registered_model_id" jsonschema:"required,description=The ID of the registered model to create the version for"`
}

func CreateModelVersionHandler(ctx context.Context, arguments CreateModelVersionHandlerArguments) (*mcpgolang.ToolResponse, error) {
	client := mr.GetInstance().GetClient(os.Getenv(defaults.ModelRegistryURL))

	createModelVersion := openapi.ModelVersionCreate{
		Name:              arguments.ModelVersionName,
		RegisteredModelId: arguments.RegisteredModelID,
	}

	modelVersion, _, err := client.ModelRegistryServiceAPI.CreateModelVersion(ctx).ModelVersionCreate(createModelVersion).Execute()
	if err != nil {
		return nil, err
	}

	outputJson, err := json.Marshal(modelVersion)
	if err != nil {
		return nil, err
	}

	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(string(outputJson))), nil
}
