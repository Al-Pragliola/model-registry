package agentcatalog

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	model "github.com/kubeflow/hub/catalog/pkg/openapi"
	"k8s.io/apimachinery/pkg/util/yaml"
)

type AgentPreviewConfig struct {
	Type           string         `json:"type" yaml:"type"`
	IncludedAgents []string       `json:"includedAgents,omitempty" yaml:"includedAgents,omitempty"`
	ExcludedAgents []string       `json:"excludedAgents,omitempty" yaml:"excludedAgents,omitempty"`
	Properties     map[string]any `json:"properties,omitempty" yaml:"properties,omitempty"`
}

func ParseAgentPreviewConfig(configBytes []byte) (*AgentPreviewConfig, error) {
	var config AgentPreviewConfig
	if err := yaml.Unmarshal(configBytes, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	if config.Type == "" {
		return nil, fmt.Errorf("missing required field: type")
	}

	if err := ValidateAgentSourceFilters(config.IncludedAgents, config.ExcludedAgents); err != nil {
		return nil, err
	}

	return &config, nil
}

func PreviewSourceAgents(_ context.Context, config *AgentPreviewConfig, catalogDataBytes []byte) ([]model.ModelPreviewResult, error) {
	agentNames, err := loadAgentNamesFromSource(config, catalogDataBytes)
	if err != nil {
		return nil, err
	}

	filter, err := NewAgentFilter(config.IncludedAgents, config.ExcludedAgents)
	if err != nil {
		return nil, fmt.Errorf("invalid filter configuration: %w", err)
	}

	results := make([]model.ModelPreviewResult, 0, len(agentNames))
	for _, name := range agentNames {
		included := filter == nil || filter.Allows(name)
		results = append(results, model.ModelPreviewResult{
			Name:     name,
			Included: included,
		})
	}

	return results, nil
}

func loadAgentNamesFromSource(config *AgentPreviewConfig, catalogDataBytes []byte) ([]string, error) {
	if config.Type != "yaml" {
		return nil, fmt.Errorf("unsupported agent source type for preview: %s", config.Type)
	}

	var catalogBytes []byte

	if len(catalogDataBytes) > 0 {
		catalogBytes = catalogDataBytes
	} else {
		path, ok := config.Properties[yamlAgentCatalogPathKey].(string)
		if !ok || path == "" {
			return nil, fmt.Errorf("missing required property: %s (provide catalogData file or set yamlCatalogPath in config)", yamlAgentCatalogPathKey)
		}

		if !filepath.IsAbs(path) {
			cwd, err := os.Getwd()
			if err != nil {
				return nil, fmt.Errorf("failed to get working directory: %w", err)
			}
			path = filepath.Join(cwd, path)
		}

		var err error
		catalogBytes, err = os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read catalog file %s: %w", path, err)
		}
	}

	var catalog yamlAgentCatalog
	if err := yaml.UnmarshalStrict(catalogBytes, &catalog); err != nil {
		return nil, fmt.Errorf("failed to parse agent catalog file: %w", err)
	}

	names := make([]string, 0, len(catalog.Agents))
	for _, a := range catalog.Agents {
		names = append(names, a.Name)
	}

	return names, nil
}
