package repositories

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/url"

	"github.com/kubeflow/hub/ui/bff/internal/integrations/httpclient"
	"github.com/kubeflow/hub/ui/bff/internal/models"
)

const agentPath = "/agents"
const agentFilterOptionPath = "/agents/filter_options"
const agentSourcesPath = "/sources"

type AgentCatalogInterface interface {
	GetAllAgents(client httpclient.HTTPClientInterface, pageValues url.Values) (*models.AgentList, error)
	GetAgentsFilter(client httpclient.HTTPClientInterface) (*models.FilterOptionsList, error)
	GetAgent(client httpclient.HTTPClientInterface, agentId string, pageValues url.Values) (*models.Agent, error)
	CreateAgentSourcePreview(client httpclient.HTTPClientInterface, payload models.CatalogSourcePreviewRequest, pageValues url.Values) (*models.CatalogSourcePreviewResult, error)
}

type AgentCatalog struct {
	AgentCatalogInterface
}

func (a *AgentCatalog) GetAllAgents(client httpclient.HTTPClientInterface, pageValues url.Values) (*models.AgentList, error) {
	responseData, err := client.GET(UrlWithPageParams(agentPath, pageValues))

	if err != nil {
		return nil, fmt.Errorf("error fetching agents list: %w", err)
	}

	var agents models.AgentList

	if err := json.Unmarshal(responseData, &agents); err != nil {
		return nil, fmt.Errorf("error decoding response data: %w", err)
	}

	return &agents, nil
}

func (a *AgentCatalog) GetAgentsFilter(client httpclient.HTTPClientInterface) (*models.FilterOptionsList, error) {
	responseData, err := client.GET(agentFilterOptionPath)

	if err != nil {
		return nil, fmt.Errorf("error fetching agent filter options: %w", err)
	}

	var filters models.FilterOptionsList

	if err := json.Unmarshal(responseData, &filters); err != nil {
		return nil, fmt.Errorf("error decoding response data: %w", err)
	}
	return &filters, nil
}

func (a *AgentCatalog) GetAgent(client httpclient.HTTPClientInterface, agentId string, pageValues url.Values) (*models.Agent, error) {
	path, err := url.JoinPath(agentPath, agentId)

	if err != nil {
		return nil, err
	}

	responseData, err := client.GET(UrlWithPageParams(path, pageValues))

	if err != nil {
		return nil, fmt.Errorf("error fetching agent: %w", err)
	}

	var agent models.Agent

	if err := json.Unmarshal(responseData, &agent); err != nil {
		return nil, fmt.Errorf("error decoding response data: %w", err)
	}

	return &agent, nil
}

func (a *AgentCatalog) CreateAgentSourcePreview(client httpclient.HTTPClientInterface, payload models.CatalogSourcePreviewRequest, pageValues url.Values) (*models.CatalogSourcePreviewResult, error) {
	path, err := url.JoinPath(agentSourcesPath, "preview")
	if err != nil {
		return nil, err
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	configPart, err := writer.CreateFormFile("config", "config.json")
	if err != nil {
		return nil, fmt.Errorf("error creating config form file: %w", err)
	}

	includedAgents := payload.IncludedAgents
	if len(includedAgents) == 0 {
		includedAgents = payload.IncludedModels
	}
	excludedAgents := payload.ExcludedAgents
	if len(excludedAgents) == 0 {
		excludedAgents = payload.ExcludedModels
	}
	configData := map[string]interface{}{
		"type":           payload.Type,
		"includedAgents": includedAgents,
		"excludedAgents": excludedAgents,
	}

	properties := make(map[string]interface{})
	var yamlContent string
	var hasYamlContent bool

	if content, ok := payload.Properties["yaml"].(string); ok && content != "" {
		yamlContent = content
		hasYamlContent = true
	} else {
		if yamlPath, ok := payload.Properties["yamlCatalogPath"].(string); ok && yamlPath != "" {
			properties["yamlCatalogPath"] = yamlPath
		} else {
			return nil, fmt.Errorf("either 'yaml' content or 'yamlCatalogPath' must be provided")
		}
	}

	if len(properties) > 0 {
		configData["properties"] = properties
	}

	configJSON, err := json.Marshal(configData)
	if err != nil {
		return nil, fmt.Errorf("error marshaling config: %w", err)
	}

	if _, err := configPart.Write(configJSON); err != nil {
		return nil, fmt.Errorf("error writing config data: %w", err)
	}

	if hasYamlContent {
		catalogDataPart, err := writer.CreateFormFile("catalogData", "catalog.yaml")
		if err != nil {
			return nil, fmt.Errorf("error creating catalogData form file: %w", err)
		}
		if _, err := catalogDataPart.Write([]byte(yamlContent)); err != nil {
			return nil, fmt.Errorf("error writing catalog data: %w", err)
		}
	}

	writer.Close()

	responseData, err := client.POSTWithContentType(UrlWithPageParams(path, pageValues), &body, writer.FormDataContentType())
	if err != nil {
		return nil, fmt.Errorf("error fetching agent source preview: %w", err)
	}

	var previewResult models.CatalogSourcePreviewResult
	if err := json.Unmarshal(responseData, &previewResult); err != nil {
		return nil, fmt.Errorf("error decoding response data: %w", err)
	}

	return &previewResult, nil
}
