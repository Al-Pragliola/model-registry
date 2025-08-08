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

type CreateModelArtifactHandlerArguments struct {
	ModelArtifactName string `json:"model_artifact_name" jsonschema:"required,description=The name of the model artifact to create"`
	URI               string `json:"uri" jsonschema:"required,description=The URI of the model artifact to create"`
	ModelVersionID    string `json:"model_version_id" jsonschema:"required,description=The ID of the model version to create the artifact for"`
}

func CreateModelArtifactHandler(ctx context.Context, arguments CreateModelArtifactHandlerArguments) (*mcpgolang.ToolResponse, error) {
	client := mr.GetInstance().GetClient(os.Getenv(defaults.ModelRegistryURL))

	// Use constructor to ensure discriminator fields (e.g., artifactType) are set
	ma := openapi.NewModelArtifact()
	ma.Name = &arguments.ModelArtifactName
	ma.Uri = &arguments.URI

	createModelArtifact := openapi.Artifact{
		ModelArtifact: ma,
	}

	modelArtifact, _, err := client.ModelRegistryServiceAPI.UpsertModelVersionArtifact(ctx, arguments.ModelVersionID).Artifact(createModelArtifact).Execute()
	if err != nil {
		return nil, err
	}

	outputJson, err := json.Marshal(modelArtifact)
	if err != nil {
		return nil, err
	}

	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(string(outputJson))), nil
}
