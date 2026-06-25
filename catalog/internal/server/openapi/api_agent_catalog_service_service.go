package openapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/kubeflow/hub/catalog/internal/catalog"
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

func (s *AgentCatalogServiceAPIService) FindAgentSources(ctx context.Context) (ImplResponse, error) {
	sources, err := s.provider.FindSources(ctx)
	if err != nil {
		return ErrorResponse(http.StatusInternalServerError, err), err
	}
	return Response(http.StatusOK, sources), nil
}

func (s *AgentCatalogServiceAPIService) PreviewAgentCatalogSource(ctx context.Context, configParam *os.File, pageSizeParam string, nextPageTokenParam string, filterStatusParam string, catalogDataParam *os.File) (ImplResponse, error) {
	pageSize := int32(10)
	if pageSizeParam != "" {
		parsed, err := strconv.ParseInt(pageSizeParam, 10, 32)
		if err != nil {
			return ErrorResponse(http.StatusBadRequest, fmt.Errorf("invalid pageSize: %w", err)), err
		}
		pageSize = int32(parsed)
	}

	filterStatus := "all"
	if filterStatusParam != "" {
		filterStatus = strings.ToLower(filterStatusParam)
		if filterStatus != "all" && filterStatus != "included" && filterStatus != "excluded" {
			err := fmt.Errorf("invalid filterStatus: must be 'all', 'included', or 'excluded'")
			return ErrorResponse(http.StatusBadRequest, err), err
		}
	}

	if configParam == nil {
		err := errors.New("config file is required")
		return ErrorResponse(http.StatusBadRequest, err), err
	}
	defer configParam.Close()

	configBytes, err := os.ReadFile(configParam.Name())
	if err != nil {
		return ErrorResponse(http.StatusBadRequest, fmt.Errorf("failed to read config file: %w", err)), err
	}

	var catalogDataBytes []byte
	if catalogDataParam != nil {
		defer catalogDataParam.Close()
		catalogDataBytes, err = os.ReadFile(catalogDataParam.Name())
		if err != nil {
			return ErrorResponse(http.StatusBadRequest, fmt.Errorf("failed to read catalogData file: %w", err)), err
		}
	}

	previewConfig, err := catalog.ParseAgentPreviewConfig(configBytes)
	if err != nil {
		return ErrorResponse(http.StatusUnprocessableEntity, fmt.Errorf("invalid config: %w", err)), err
	}

	previewResults, err := catalog.PreviewSourceAgents(ctx, previewConfig, catalogDataBytes)
	if err != nil {
		return ErrorResponse(http.StatusUnprocessableEntity, fmt.Errorf("failed to load agents: %w", err)), err
	}

	filteredResults := make([]model.ModelPreviewResult, 0)
	for _, result := range previewResults {
		item := model.ModelPreviewResult{Name: result.Name, Included: result.Included}
		switch filterStatus {
		case "included":
			if result.Included {
				filteredResults = append(filteredResults, item)
			}
		case "excluded":
			if !result.Included {
				filteredResults = append(filteredResults, item)
			}
		default:
			filteredResults = append(filteredResults, item)
		}
	}

	var includedCount, excludedCount int32
	for _, result := range previewResults {
		if result.Included {
			includedCount++
		} else {
			excludedCount++
		}
	}
	totalCount := int32(len(previewResults))

	start := int32(0)
	if nextPageTokenParam != "" {
		parsed, err := strconv.ParseInt(nextPageTokenParam, 10, 32)
		if err == nil {
			start = int32(parsed)
		}
	}

	end := start + pageSize
	if end > int32(len(filteredResults)) {
		end = int32(len(filteredResults))
	}
	if start > int32(len(filteredResults)) {
		start = int32(len(filteredResults))
	}

	pagedResults := filteredResults[start:end]

	var nextToken string
	if end < int32(len(filteredResults)) {
		nextToken = fmt.Sprintf("%d", end)
	}

	response := model.CatalogSourcePreviewResponse{
		Items: pagedResults,
		Summary: model.CatalogSourcePreviewResponseAllOfSummary{
			TotalModels:    totalCount,
			IncludedModels: includedCount,
			ExcludedModels: excludedCount,
		},
		NextPageToken: nextToken,
		PageSize:      pageSize,
		Size:          int32(len(pagedResults)),
	}

	return Response(http.StatusOK, response), nil
}
