package openapi

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	model "github.com/kubeflow/hub/catalog/pkg/openapi"
)

var fixedCSVColumns = []string{
	"id",
	"name",
	"source_id",
	"description",
	"provider",
	"license",
	"license_link",
	"library_name",
	"maturity",
	"language",
	"tasks",
	"create_time",
	"last_update_time",
}

func collectCustomPropertyKeys(models []model.CatalogModel) []string {
	keySet := make(map[string]struct{})
	for i := range models {
		for key := range models[i].CustomProperties {
			keySet[key] = struct{}{}
		}
	}
	keys := make([]string, 0, len(keySet))
	for key := range keySet {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func buildCSVHeaders(customKeys []string) []string {
	headers := make([]string, 0, len(fixedCSVColumns)+len(customKeys))
	headers = append(headers, fixedCSVColumns...)
	for _, key := range customKeys {
		headers = append(headers, "custom_"+key)
	}
	return headers
}

func metadataValueToString(v model.MetadataValue) string {
	if v.MetadataStringValue != nil {
		return v.MetadataStringValue.StringValue
	}
	if v.MetadataIntValue != nil {
		return v.MetadataIntValue.IntValue
	}
	if v.MetadataDoubleValue != nil {
		return fmt.Sprintf("%g", v.MetadataDoubleValue.DoubleValue)
	}
	if v.MetadataBoolValue != nil {
		return fmt.Sprintf("%t", v.MetadataBoolValue.BoolValue)
	}
	if v.MetadataStructValue != nil {
		return v.MetadataStructValue.StructValue
	}
	if v.MetadataProtoValue != nil {
		data, err := json.Marshal(v.MetadataProtoValue)
		if err != nil {
			return ""
		}
		return string(data)
	}
	return ""
}

func derefString(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func modelToCSVRow(m model.CatalogModel, customKeys []string) []string {
	row := []string{
		derefString(m.Id),
		m.Name,
		derefString(m.SourceId),
		derefString(m.Description),
		derefString(m.Provider),
		derefString(m.License),
		derefString(m.LicenseLink),
		derefString(m.LibraryName),
		derefString(m.Maturity),
		strings.Join(m.Language, ";"),
		strings.Join(m.Tasks, ";"),
		derefString(m.CreateTimeSinceEpoch),
		derefString(m.LastUpdateTimeSinceEpoch),
	}
	for _, key := range customKeys {
		val, ok := m.CustomProperties[key]
		if !ok {
			row = append(row, "")
		} else {
			row = append(row, metadataValueToString(val))
		}
	}
	return row
}

func writeCSV(w io.Writer, models []model.CatalogModel) error {
	// UTF-8 BOM for Excel compatibility
	if _, err := w.Write([]byte{0xEF, 0xBB, 0xBF}); err != nil {
		return err
	}

	customKeys := collectCustomPropertyKeys(models)
	headers := buildCSVHeaders(customKeys)

	csvWriter := csv.NewWriter(w)
	if err := csvWriter.Write(headers); err != nil {
		return err
	}
	for i := range models {
		if err := csvWriter.Write(modelToCSVRow(models[i], customKeys)); err != nil {
			return err
		}
	}
	csvWriter.Flush()
	return csvWriter.Error()
}

func writeCSVToFile(models []model.CatalogModel) (*os.File, error) {
	f, err := os.CreateTemp("", "model-catalog-export-*.csv")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}

	if err := writeCSV(f, models); err != nil {
		_ = f.Close()
		_ = os.Remove(f.Name())
		return nil, fmt.Errorf("failed to write CSV: %w", err)
	}

	// Seek back to start so EncodeJSONResponse can read the file
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		_ = f.Close()
		_ = os.Remove(f.Name())
		return nil, fmt.Errorf("failed to seek: %w", err)
	}

	return f, nil
}
