package openapi

import (
	"context"
	"errors"
	"net/http"

	agentcatalog "github.com/kubeflow/hub/catalog/internal/catalog/agentcatalog"
	model "github.com/kubeflow/hub/catalog/pkg/openapi"
	"github.com/kubeflow/hub/pkg/api"
)

type AgentCatalogServiceAPIService struct {
	provider     *agentcatalog.DBAgentCatalog
	agentSources *agentcatalog.AgentSourceCollection
}

var _ AgentCatalogServiceAPIServicer = &AgentCatalogServiceAPIService{}

func NewAgentCatalogServiceAPIService(provider *agentcatalog.DBAgentCatalog, agentSources *agentcatalog.AgentSourceCollection) AgentCatalogServiceAPIServicer {
	return &AgentCatalogServiceAPIService{
		provider:     provider,
		agentSources: agentSources,
	}
}

func (s *AgentCatalogServiceAPIService) FindAgents(ctx context.Context, name string, q string, sourceLabel []string, filterQuery string, pageSize string, orderBy model.OrderByField, sortOrder model.SortOrder, nextPageToken string) (ImplResponse, error) {
	pageSizeInt, err := parsePaginationParams(pageSize, nextPageToken)
	if err != nil {
		return ErrorResponse(http.StatusBadRequest, err), err
	}

	if len(sourceLabel) == 1 && sourceLabel[0] == "" {
		sourceLabel = nil
	}

	var sourceIDs []string
	if len(sourceLabel) > 0 && s.agentSources != nil {
		sources := s.agentSources.ByLabel(sourceLabel)
		if len(sources) == 0 {
			return Response(http.StatusOK, model.AgentList{
				Items:    []model.Agent{},
				PageSize: &pageSizeInt,
			}), nil
		}
		sourceIDs = make([]string, len(sources))
		for i, source := range sources {
			sourceIDs[i] = source.ID
		}
	}

	params := agentcatalog.ListAgentsParams{
		Name:          name,
		Query:         q,
		SourceIDs:     sourceIDs,
		FilterQuery:   filterQuery,
		PageSize:      pageSizeInt,
		OrderBy:       orderBy,
		SortOrder:     sortOrder,
		NextPageToken: &nextPageToken,
	}

	agents, err := s.provider.ListAgents(ctx, params)
	if err != nil {
		return ErrorResponse(api.ErrToStatus(err), err), err
	}

	return Response(http.StatusOK, agents), nil
}

func (s *AgentCatalogServiceAPIService) GetAgent(ctx context.Context, id string) (ImplResponse, error) {
	agent, err := s.provider.GetAgent(ctx, id)
	if err != nil {
		if errors.Is(err, api.ErrNotFound) {
			return ErrorResponse(http.StatusNotFound, err), err
		}
		return ErrorResponse(api.ErrToStatus(err), err), err
	}

	if agent == nil {
		return ErrorResponse(http.StatusNotFound, errors.New("agent not found")), nil
	}

	return Response(http.StatusOK, agent), nil
}

func (s *AgentCatalogServiceAPIService) GetAgentFilterOptions(ctx context.Context) (ImplResponse, error) {
	filterOptions, err := s.provider.GetFilterOptions(ctx)
	if err != nil {
		return ErrorResponse(http.StatusInternalServerError, err), err
	}
	return Response(http.StatusOK, *filterOptions), nil
}
