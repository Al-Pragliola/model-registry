package api

import (
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/kubeflow/hub/ui/bff/internal/constants"
	"github.com/kubeflow/hub/ui/bff/internal/integrations/httpclient"
	"github.com/kubeflow/hub/ui/bff/internal/models"
	"github.com/kubeflow/model-registry/pkg/openapi"
)

var fixedCSVColumns = []string{
	"Name",
	"Source ID",
	"Provider",
	"Description",
	"Maturity",
	"License",
	"License Link",
	"Library",
	"Language",
	"Tasks",
	"Created",
	"Last Updated",
}

func (app *App) ExportCatalogModelsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	client, ok := r.Context().Value(constants.ModelCatalogHttpClientKey).(httpclient.HTTPClientInterface)
	if !ok {
		app.serverErrorResponse(w, r, errors.New("catalog REST client not found"))
		return
	}

	allModels, err := app.repositories.ModelCatalogClient.GetAllCatalogModelsForExport(client, r.URL.Query())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	customPropKeys := collectCustomPropertyKeys(allModels)

	headers := make([]string, 0, len(fixedCSVColumns)+len(customPropKeys))
	headers = append(headers, fixedCSVColumns...)
	headers = append(headers, customPropKeys...)

	filename := fmt.Sprintf("model-catalog-export-%s.csv", time.Now().UTC().Format("2006-01-02"))
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filename))

	// UTF-8 BOM for Excel compatibility
	_, _ = w.Write([]byte{0xEF, 0xBB, 0xBF})

	csvWriter := csv.NewWriter(w)

	if err := csvWriter.Write(headers); err != nil {
		app.logger.Error("failed to write CSV header", "error", err)
		return
	}

	for _, model := range allModels {
		row := modelToCSVRow(model, customPropKeys)
		if err := csvWriter.Write(row); err != nil {
			app.logger.Error("failed to write CSV row", "error", err, "model", model.Name)
			return
		}
	}

	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		app.logger.Error("CSV flush error", "error", err)
	}
}

func collectCustomPropertyKeys(models []models.CatalogModel) []string {
	keySet := make(map[string]struct{})
	for _, m := range models {
		if m.CustomProperties == nil {
			continue
		}
		for key := range *m.CustomProperties {
			keySet[key] = struct{}{}
		}
	}

	keys := make([]string, 0, len(keySet))
	for k := range keySet {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func modelToCSVRow(model models.CatalogModel, customPropKeys []string) []string {
	row := make([]string, 0, len(fixedCSVColumns)+len(customPropKeys))

	row = append(row, model.Name)
	row = append(row, derefStr(model.SourceId))
	row = append(row, derefStr(model.Provider))
	row = append(row, derefStr(model.Description))
	row = append(row, derefStr(model.Maturity))
	row = append(row, derefStr(model.License))
	row = append(row, derefStr(model.LicenseLink))
	row = append(row, derefStr(model.LibraryName))
	row = append(row, strings.Join(model.Language, "; "))
	row = append(row, strings.Join(model.Tasks, "; "))
	row = append(row, epochMsToISO(model.CreateTimeSinceEpoch))
	row = append(row, epochMsToISO(model.LastUpdateTimeSinceEpoch))

	for _, key := range customPropKeys {
		row = append(row, extractCustomPropertyValue(model.CustomProperties, key))
	}

	return row
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func epochMsToISO(epoch *string) string {
	if epoch == nil || *epoch == "" {
		return ""
	}
	var ms int64
	if _, err := fmt.Sscanf(*epoch, "%d", &ms); err != nil {
		return *epoch
	}
	return time.UnixMilli(ms).UTC().Format(time.RFC3339)
}

func extractCustomPropertyValue(props *map[string]openapi.MetadataValue, key string) string {
	if props == nil {
		return ""
	}
	val, ok := (*props)[key]
	if !ok {
		return ""
	}

	if val.MetadataStringValue != nil {
		return val.MetadataStringValue.StringValue
	}
	if val.MetadataIntValue != nil {
		return val.MetadataIntValue.IntValue
	}
	if val.MetadataDoubleValue != nil {
		return fmt.Sprintf("%g", val.MetadataDoubleValue.DoubleValue)
	}
	if val.MetadataBoolValue != nil {
		return fmt.Sprintf("%t", val.MetadataBoolValue.BoolValue)
	}

	return ""
}
