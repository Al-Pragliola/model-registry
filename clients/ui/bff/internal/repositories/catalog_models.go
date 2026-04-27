package repositories

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/kubeflow/hub/ui/bff/internal/integrations/httpclient"
	"github.com/kubeflow/hub/ui/bff/internal/models"
)

const catalogModelsPath = "/models"

const exportPageSize = 100

type CatalogModelsInterface interface {
	GetAllCatalogModelsAcrossSources(client httpclient.HTTPClientInterface, pageValues url.Values) (*models.CatalogModelList, error)
	GetAllCatalogModelsForExport(client httpclient.HTTPClientInterface, pageValues url.Values) ([]models.CatalogModel, error)
}

type CatalogModels struct {
	CatalogModelsInterface
}

func (a CatalogModels) GetAllCatalogModelsAcrossSources(client httpclient.HTTPClientInterface, pageValues url.Values) (*models.CatalogModelList, error) {
	responseData, err := client.GET(UrlWithPageParams(catalogModelsPath, pageValues))
	if err != nil {
		return nil, fmt.Errorf("error fetching sourcesPath: %w", err)
	}

	var models models.CatalogModelList

	if err := json.Unmarshal(responseData, &models); err != nil {
		return nil, fmt.Errorf("error decoding response data: %w", err)
	}

	return &models, nil
}

func (a CatalogModels) GetAllCatalogModelsForExport(client httpclient.HTTPClientInterface, pageValues url.Values) ([]models.CatalogModel, error) {
	exportValues := FilterPageValues(pageValues)
	exportValues.Set("pageSize", fmt.Sprintf("%d", exportPageSize))
	exportValues.Del("nextPageToken")

	var allModels []models.CatalogModel

	for {
		responseData, err := client.GET(UrlWithParams(catalogModelsPath, exportValues))
		if err != nil {
			return nil, fmt.Errorf("error fetching models for export: %w", err)
		}

		var page models.CatalogModelList
		if err := json.Unmarshal(responseData, &page); err != nil {
			return nil, fmt.Errorf("error decoding response data: %w", err)
		}

		allModels = append(allModels, page.Items...)

		if page.NextPageToken == "" {
			break
		}
		exportValues.Set("nextPageToken", page.NextPageToken)
	}

	return allModels, nil
}
