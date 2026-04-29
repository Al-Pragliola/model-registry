package repositories

import (
	"fmt"
	"net/url"

	"github.com/kubeflow/hub/ui/bff/internal/integrations/httpclient"
)

const catalogExportPath = "/models/export"

type CatalogExportInterface interface {
	ExportCatalogModels(client httpclient.HTTPClientInterface, pageValues url.Values) ([]byte, error)
}

type CatalogExport struct {
	CatalogExportInterface
}

func (a CatalogExport) ExportCatalogModels(client httpclient.HTTPClientInterface, pageValues url.Values) ([]byte, error) {
	responseData, err := client.GET(UrlWithPageParams(catalogExportPath, pageValues))
	if err != nil {
		return nil, fmt.Errorf("error exporting catalog models: %w", err)
	}
	return responseData, nil
}
