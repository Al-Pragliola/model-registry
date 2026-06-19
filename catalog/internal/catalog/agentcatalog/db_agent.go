package agentcatalog

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/kubeflow/hub/catalog/internal/catalog/agentcatalog/models"
	openapi "github.com/kubeflow/hub/catalog/pkg/openapi"
	"github.com/kubeflow/hub/internal/platform/apiutils"
	"github.com/kubeflow/hub/pkg/api"
)

type DBAgentCatalog struct {
	services Services
	sources  *AgentSourceCollection
}

func NewDBAgentCatalog(services Services, sources *AgentSourceCollection) *DBAgentCatalog {
	return &DBAgentCatalog{
		services: services,
		sources:  sources,
	}
}

type ListAgentsParams struct {
	Name          string
	Query         string
	SourceIDs     []string
	FilterQuery   string
	OrderBy       openapi.OrderByField
	SortOrder     openapi.SortOrder
	NextPageToken *string
	PageSize      int32
}

func (d *DBAgentCatalog) ListAgents(_ context.Context, params ListAgentsParams) (openapi.AgentList, error) {
	listOptions := models.AgentListOptions{
		FilterQuery: &params.FilterQuery,
	}

	if params.Name != "" {
		listOptions.Name = &params.Name
	}

	if params.Query != "" {
		listOptions.Query = &params.Query
	}

	if len(params.SourceIDs) > 0 {
		listOptions.SourceIDs = &params.SourceIDs
	}

	orderBy := strings.ToUpper(string(params.OrderBy))
	sortOrder := strings.ToUpper(string(params.SortOrder))
	listOptions.Pagination.PageSize = &params.PageSize
	if orderBy != "" {
		listOptions.Pagination.OrderBy = &orderBy
	}
	if sortOrder != "" {
		listOptions.Pagination.SortOrder = &sortOrder
	}
	if params.NextPageToken != nil {
		listOptions.Pagination.NextPageToken = params.NextPageToken
	}

	result, err := d.services.AgentRepository.List(&listOptions)
	if err != nil {
		return openapi.AgentList{}, err
	}

	apiAgents := make([]openapi.Agent, 0)
	for _, dbAgent := range result.Items {
		apiAgent := mapDBAgentToAPI(dbAgent)
		apiAgents = append(apiAgents, apiAgent)
	}

	size := int32(len(apiAgents))
	agentList := openapi.AgentList{
		Items:    apiAgents,
		Size:     &size,
		PageSize: &params.PageSize,
	}
	if result.NextPageToken != "" {
		agentList.NextPageToken = &result.NextPageToken
	}

	return agentList, nil
}

func (d *DBAgentCatalog) GetAgent(_ context.Context, agentID string) (*openapi.Agent, error) {
	id, err := apiutils.ValidateIDAsInt32(agentID, "agent")
	if err != nil {
		return nil, fmt.Errorf("invalid agent ID '%s': %w", agentID, api.ErrBadRequest)
	}

	dbAgent, err := d.services.AgentRepository.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("agent not found with ID %s: %w", agentID, api.ErrNotFound)
	}

	apiAgent := mapDBAgentToAPI(dbAgent)
	return &apiAgent, nil
}

func (d *DBAgentCatalog) GetFilterOptions(_ context.Context) (*openapi.FilterOptionsList, error) {
	options := make(map[string]openapi.FilterOption)
	return &openapi.FilterOptionsList{
		Filters: &options,
	}, nil
}

func (d *DBAgentCatalog) FindSources(_ context.Context) (openapi.CatalogSourceList, error) {
	allSources := d.sources.AllSources()
	sources := make([]openapi.CatalogSource, 0)

	for id, source := range allSources {
		catalogSource := openapi.CatalogSource{
			Id:   id,
			Name: source.Name,
		}
		if len(source.Labels) > 0 {
			catalogSource.Labels = source.Labels
		}

		sources = append(sources, catalogSource)
	}

	return openapi.CatalogSourceList{
		Items: sources,
		Size:  int32(len(sources)),
	}, nil
}

func mapDBAgentToAPI(dbAgent models.Agent) openapi.Agent {
	res := openapi.Agent{}

	if dbAgent.GetID() != nil {
		idStr := fmt.Sprintf("%d", *dbAgent.GetID())
		res.Id = &idStr
	}

	attrs := dbAgent.GetAttributes()
	if attrs != nil {
		if attrs.Name != nil {
			res.Name = *attrs.Name
		}
		if attrs.ExternalID != nil {
			res.ExternalId = attrs.ExternalID
		}
		if attrs.CreateTimeSinceEpoch != nil {
			ts := fmt.Sprintf("%d", *attrs.CreateTimeSinceEpoch)
			res.CreateTimeSinceEpoch = &ts
		}
		if attrs.LastUpdateTimeSinceEpoch != nil {
			ts := fmt.Sprintf("%d", *attrs.LastUpdateTimeSinceEpoch)
			res.LastUpdateTimeSinceEpoch = &ts
		}
	}

	if dbAgent.GetProperties() != nil {
		for _, prop := range *dbAgent.GetProperties() {
			if prop.StringValue == nil {
				continue
			}
			switch prop.Name {
			case "source_id":
				res.SourceId = prop.StringValue
			case "description":
				res.Description = prop.StringValue
			case "displayName":
				res.DisplayName = prop.StringValue
			case "framework":
				res.Framework = prop.StringValue
			case "agentType":
				res.AgentType = prop.StringValue
			case "logo":
				res.Logo = prop.StringValue
			case "repositoryUrl":
				res.RepositoryUrl = prop.StringValue
			case "publishedDate":
				if t, err := time.Parse(time.RFC3339, *prop.StringValue); err == nil {
					res.PublishedDate = &t
				}
			case "readme":
				res.Readme = prop.StringValue
			case "tags":
				var tags []string
				if err := json.Unmarshal([]byte(*prop.StringValue), &tags); err == nil {
					res.Tags = tags
				}
			case "models":
				var models []string
				if err := json.Unmarshal([]byte(*prop.StringValue), &models); err == nil {
					res.Models = models
				}
			case "env":
				var env []openapi.AgentEnvVar
				if err := json.Unmarshal([]byte(*prop.StringValue), &env); err == nil {
					res.Env = env
				}
			case "artifacts":
				var artifacts []openapi.AgentArtifact
				if err := json.Unmarshal([]byte(*prop.StringValue), &artifacts); err == nil {
					res.Artifacts = artifacts
				}
			}
		}
	}

	return res
}
