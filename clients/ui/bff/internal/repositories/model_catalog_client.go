package repositories

import (
	"log/slog"
)

type ModelCatalogClientInterface interface {
	CatalogSourcesInterface
	CatalogModelsInterface
	CatalogExportInterface
	CatalogSourcePreviewInterface
	McpServerCatalogInterface
}

type ModelCatalogClient struct {
	logger *slog.Logger
	CatalogSources
	CatalogModels
	CatalogExport
	CatalogSourcePreview
	McpServerCatalog
}

func NewModelCatalogClient(logger *slog.Logger) (ModelCatalogClientInterface, error) {
	return &ModelCatalogClient{logger: logger}, nil
}
