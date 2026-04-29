package openapi

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/kubeflow/hub/catalog/internal/catalog"
	model "github.com/kubeflow/hub/catalog/pkg/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newExportTestService(mockModels map[string]*model.CatalogModel) *ModelCatalogServiceAPIService {
	provider := &mockModelProvider{models: mockModels}
	sources := catalog.NewSourceCollection()
	sources.Merge("", map[string]catalog.ModelSource{
		"source1": {CatalogSource: model.CatalogSource{Id: "source1", Name: "Source 1"}},
	})
	labels := catalog.NewLabelCollection()
	svc := NewModelCatalogServiceAPIService(provider, sources, nil, labels, nil)
	return svc.(*ModelCatalogServiceAPIService)
}

func TestExportModelsCSV(t *testing.T) {
	id1 := "1"
	id2 := "2"
	models := map[string]*model.CatalogModel{
		"Model A": {Id: &id1, Name: "Model A", Provider: strPtr("IBM"), Language: []string{"en"}},
		"Model B": {Id: &id2, Name: "Model B", Provider: strPtr("Meta"), Tasks: []string{"text-gen"}},
	}

	svc := newExportTestService(models)
	resp, err := svc.ExportModels(context.Background(), false, nil, nil, nil, "", nil, "", "", model.ORDERBYFIELD_NAME, model.SORTORDER_ASC, "")

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)

	f, ok := resp.Body.(*os.File)
	require.True(t, ok, "response body should be *os.File")
	defer func() { _ = f.Close() }()
	defer func() { _ = os.Remove(f.Name()) }()
}

func TestExportModelsDryRun(t *testing.T) {
	id1 := "1"
	id2 := "2"
	models := map[string]*model.CatalogModel{
		"Model A": {
			Id: &id1, Name: "Model A",
			CustomProperties: map[string]model.MetadataValue{
				"arch": model.MetadataStringValueAsMetadataValue(model.NewMetadataStringValue("transformer", "MetadataStringValue")),
			},
		},
		"Model B": {Id: &id2, Name: "Model B"},
	}

	svc := newExportTestService(models)
	resp, err := svc.ExportModels(context.Background(), true, nil, nil, nil, "", nil, "", "10", model.ORDERBYFIELD_NAME, model.SORTORDER_ASC, "")

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)

	dryRun, ok := resp.Body.(model.ExportDryRunResponse)
	require.True(t, ok, "response body should be ExportDryRunResponse")

	assert.Equal(t, int32(2), dryRun.TotalCount)
	assert.Equal(t, int32(2), dryRun.Size)
	assert.Contains(t, dryRun.Columns, "custom_arch")
	assert.Contains(t, dryRun.Columns, "name")
}

func TestExportModelsFilterByIDs(t *testing.T) {
	id1 := "1"
	id2 := "2"
	id3 := "3"
	models := map[string]*model.CatalogModel{
		"Model A": {Id: &id1, Name: "Model A"},
		"Model B": {Id: &id2, Name: "Model B"},
		"Model C": {Id: &id3, Name: "Model C"},
	}

	svc := newExportTestService(models)
	resp, err := svc.ExportModels(context.Background(), true, []string{"1", "3"}, nil, nil, "", nil, "", "10", model.ORDERBYFIELD_NAME, model.SORTORDER_ASC, "")

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)

	dryRun := resp.Body.(model.ExportDryRunResponse)
	assert.Equal(t, int32(2), dryRun.TotalCount)
	require.Len(t, dryRun.Items, 2)

	names := make([]string, len(dryRun.Items))
	for i, item := range dryRun.Items {
		names[i] = item.Name
	}
	assert.Contains(t, names, "Model A")
	assert.Contains(t, names, "Model C")
}

func TestExportModelsExcludeIDs(t *testing.T) {
	id1 := "1"
	id2 := "2"
	id3 := "3"
	models := map[string]*model.CatalogModel{
		"Model A": {Id: &id1, Name: "Model A"},
		"Model B": {Id: &id2, Name: "Model B"},
		"Model C": {Id: &id3, Name: "Model C"},
	}

	svc := newExportTestService(models)
	resp, err := svc.ExportModels(context.Background(), true, nil, []string{"2"}, nil, "", nil, "", "10", model.ORDERBYFIELD_NAME, model.SORTORDER_ASC, "")

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)

	dryRun := resp.Body.(model.ExportDryRunResponse)
	assert.Equal(t, int32(2), dryRun.TotalCount)
	require.Len(t, dryRun.Items, 2)

	for _, item := range dryRun.Items {
		assert.NotEqual(t, "Model B", item.Name)
	}
}

func TestExportModelsMutualExclusivityIDAndExcludeID(t *testing.T) {
	svc := newExportTestService(nil)
	resp, err := svc.ExportModels(context.Background(), false, []string{"1"}, []string{"2"}, nil, "", nil, "", "", "", "", "")

	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestExportModelsMutualExclusivityIDAndFilters(t *testing.T) {
	svc := newExportTestService(nil)
	resp, err := svc.ExportModels(context.Background(), false, []string{"1"}, nil, []string{"source1"}, "", nil, "", "", "", "", "")

	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestExportModelsEmptyResult(t *testing.T) {
	svc := newExportTestService(map[string]*model.CatalogModel{})

	// CSV mode
	resp, err := svc.ExportModels(context.Background(), false, nil, nil, nil, "", nil, "", "", model.ORDERBYFIELD_NAME, model.SORTORDER_ASC, "")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)
	f, ok := resp.Body.(*os.File)
	require.True(t, ok)
	defer func() { _ = f.Close() }()
	defer func() { _ = os.Remove(f.Name()) }()

	// Dry run mode
	resp, err = svc.ExportModels(context.Background(), true, nil, nil, nil, "", nil, "", "10", model.ORDERBYFIELD_NAME, model.SORTORDER_ASC, "")
	require.NoError(t, err)
	dryRun := resp.Body.(model.ExportDryRunResponse)
	assert.Equal(t, int32(0), dryRun.TotalCount)
	assert.Empty(t, dryRun.Items)
}

func TestExportModelsDryRunPagination(t *testing.T) {
	models := make(map[string]*model.CatalogModel)
	for i := range 5 {
		id := strPtr(string(rune('1' + i)))
		name := string(rune('A' + i))
		models[name] = &model.CatalogModel{Id: id, Name: name}
	}

	svc := newExportTestService(models)

	// First page: 2 items
	resp, err := svc.ExportModels(context.Background(), true, nil, nil, nil, "", nil, "", "2", model.ORDERBYFIELD_NAME, model.SORTORDER_ASC, "")
	require.NoError(t, err)

	dryRun := resp.Body.(model.ExportDryRunResponse)
	assert.Equal(t, int32(5), dryRun.TotalCount)
	assert.Equal(t, int32(2), dryRun.Size)
	assert.NotEmpty(t, dryRun.NextPageToken)
}
