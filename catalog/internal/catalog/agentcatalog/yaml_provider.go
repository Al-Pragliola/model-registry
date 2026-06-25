package agentcatalog

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang/glog"
	"github.com/kubeflow/hub/catalog/internal/catalog/agentcatalog/models"
	"github.com/kubeflow/hub/catalog/internal/catalog/basecatalog"
	mrmodels "github.com/kubeflow/hub/internal/platform/db/entity"
	"k8s.io/apimachinery/pkg/util/yaml"
)

const YamlAgentCatalogPathKey = "yamlCatalogPath"

type YamlAgentEnvVar struct {
	Name     string `yaml:"name" json:"name"`
	Required bool   `yaml:"required" json:"required"`
}

type YamlAgentArtifact struct {
	URI string `yaml:"uri" json:"uri"`
}

type YamlAgent struct {
	Name                     string             `yaml:"name" json:"name"`
	DisplayName              *string            `yaml:"displayName,omitempty" json:"displayName,omitempty"`
	Description              *string            `yaml:"description,omitempty" json:"description,omitempty"`
	Readme                   *string            `yaml:"readme,omitempty" json:"readme,omitempty"`
	Framework                *string            `yaml:"framework,omitempty" json:"framework,omitempty"`
	AgentType                *string            `yaml:"agentType,omitempty" json:"agentType,omitempty"`
	Tags                     []string           `yaml:"tags,omitempty" json:"tags,omitempty"`
	Models                   []string           `yaml:"models,omitempty" json:"models,omitempty"`
	Logo                     *string            `yaml:"logo,omitempty" json:"logo,omitempty"`
	RepositoryUrl            *string            `yaml:"repositoryUrl,omitempty" json:"repositoryUrl,omitempty"`
	PublishedDate            *string            `yaml:"publishedDate,omitempty" json:"publishedDate,omitempty"`
	Env                      []YamlAgentEnvVar  `yaml:"env,omitempty" json:"env,omitempty"`
	Artifacts                []YamlAgentArtifact `yaml:"artifacts,omitempty" json:"artifacts,omitempty"`
	CreateTimeSinceEpoch     *string            `yaml:"createTimeSinceEpoch,omitempty" json:"createTimeSinceEpoch,omitempty"`
	LastUpdateTimeSinceEpoch *string            `yaml:"lastUpdateTimeSinceEpoch,omitempty" json:"lastUpdateTimeSinceEpoch,omitempty"`
}

type YamlAgentCatalog struct {
	Source string      `yaml:"source" json:"source"`
	Agents []YamlAgent `yaml:"agents" json:"agents"`
}

func (l *AgentLoader) loadFromYAML(_ context.Context, sourceID string, source basecatalog.AgentSource) error {
	yamlPath, ok := source.Properties[YamlAgentCatalogPathKey].(string)
	if !ok {
		return fmt.Errorf("yamlCatalogPath property is required for YAML agent provider")
	}

	if !filepath.IsAbs(yamlPath) {
		sourceDir := filepath.Dir(source.Origin)
		yamlPath = filepath.Join(sourceDir, yamlPath)
	}

	data, err := os.ReadFile(yamlPath)
	if err != nil {
		return fmt.Errorf("error reading YAML file %s: %w", yamlPath, err)
	}

	var catalog YamlAgentCatalog
	if err := yaml.Unmarshal(data, &catalog); err != nil {
		return fmt.Errorf("error parsing YAML from %s: %w", yamlPath, err)
	}

	glog.Infof("loading %d agents from source %s", len(catalog.Agents), sourceID)

	filter, err := NewAgentFilter(source.IncludedAgents, source.ExcludedAgents)
	if err != nil {
		return fmt.Errorf("invalid agent filter for source %s: %w", sourceID, err)
	}

	l.services.AgentRepository.DeleteBySource(sourceID)

	for _, ya := range catalog.Agents {
		if filter != nil && !filter.Allows(ya.Name) {
			glog.V(2).Infof("agent %s excluded by filter in source %s", ya.Name, sourceID)
			continue
		}

		agent := convertYAMLAgentToEntity(ya, sourceID)

		_, err := l.services.AgentRepository.Save(agent)
		if err != nil {
			glog.Errorf("error saving agent %s from source %s: %v", ya.Name, sourceID, err)
			continue
		}
	}

	glog.Infof("loaded agents from source %s", sourceID)
	return nil
}

func convertYAMLAgentToEntity(ya YamlAgent, sourceID string) models.Agent {
	attrs := &models.AgentAttributes{
		Name: &ya.Name,
	}

	if ya.CreateTimeSinceEpoch != nil {
		if ts, err := parseInt64(*ya.CreateTimeSinceEpoch); err == nil {
			attrs.CreateTimeSinceEpoch = &ts
		}
	}
	if ya.LastUpdateTimeSinceEpoch != nil {
		if ts, err := parseInt64(*ya.LastUpdateTimeSinceEpoch); err == nil {
			attrs.LastUpdateTimeSinceEpoch = &ts
		}
	}

	agent := &models.AgentImpl{
		Attributes: attrs,
	}

	properties := []mrmodels.Properties{
		mrmodels.NewStringProperty("source_id", sourceID, false),
	}

	if ya.Description != nil {
		properties = append(properties, mrmodels.NewStringProperty("description", *ya.Description, false))
	}
	if ya.DisplayName != nil {
		properties = append(properties, mrmodels.NewStringProperty("displayName", *ya.DisplayName, false))
	}
	if ya.Framework != nil {
		properties = append(properties, mrmodels.NewStringProperty("framework", *ya.Framework, false))
	}
	if ya.AgentType != nil {
		properties = append(properties, mrmodels.NewStringProperty("agentType", *ya.AgentType, false))
	}
	if ya.Logo != nil {
		properties = append(properties, mrmodels.NewStringProperty("logo", *ya.Logo, false))
	}
	if ya.RepositoryUrl != nil {
		properties = append(properties, mrmodels.NewStringProperty("repositoryUrl", *ya.RepositoryUrl, false))
	}
	if ya.PublishedDate != nil {
		properties = append(properties, mrmodels.NewStringProperty("publishedDate", *ya.PublishedDate, false))
	}
	if ya.Readme != nil {
		properties = append(properties, mrmodels.NewStringProperty("readme", *ya.Readme, false))
	}

	if len(ya.Tags) > 0 {
		if jsonBytes, err := json.Marshal(ya.Tags); err == nil {
			properties = append(properties, mrmodels.NewStringProperty("tags", string(jsonBytes), false))
		}
	}
	if len(ya.Models) > 0 {
		if jsonBytes, err := json.Marshal(ya.Models); err == nil {
			properties = append(properties, mrmodels.NewStringProperty("models", string(jsonBytes), false))
		}
	}
	if len(ya.Env) > 0 {
		if jsonBytes, err := json.Marshal(ya.Env); err == nil {
			properties = append(properties, mrmodels.NewStringProperty("env", string(jsonBytes), false))
		}
	}
	if len(ya.Artifacts) > 0 {
		if jsonBytes, err := json.Marshal(ya.Artifacts); err == nil {
			properties = append(properties, mrmodels.NewStringProperty("artifacts", string(jsonBytes), false))
		}
	}

	agent.Properties = &properties
	return agent
}

func parseInt64(s string) (int64, error) {
	var v int64
	_, err := fmt.Sscanf(s, "%d", &v)
	return v, err
}
