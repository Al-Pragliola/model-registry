package openapi

import (
	"bytes"
	"encoding/csv"
	"os"
	"strings"
	"testing"

	model "github.com/kubeflow/hub/catalog/pkg/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func strPtr(s string) *string { return &s }

func TestCollectCustomPropertyKeys(t *testing.T) {
	models := []model.CatalogModel{
		{
			Name: "model-a",
			CustomProperties: map[string]model.MetadataValue{
				"arch":  model.MetadataStringValueAsMetadataValue(model.NewMetadataStringValue("transformer", "MetadataStringValue")),
				"zebra": model.MetadataStringValueAsMetadataValue(model.NewMetadataStringValue("z", "MetadataStringValue")),
			},
		},
		{
			Name: "model-b",
			CustomProperties: map[string]model.MetadataValue{
				"arch":  model.MetadataStringValueAsMetadataValue(model.NewMetadataStringValue("cnn", "MetadataStringValue")),
				"alpha": model.MetadataStringValueAsMetadataValue(model.NewMetadataStringValue("a", "MetadataStringValue")),
			},
		},
	}

	keys := collectCustomPropertyKeys(models)
	assert.Equal(t, []string{"alpha", "arch", "zebra"}, keys)
}

func TestCollectCustomPropertyKeysEmpty(t *testing.T) {
	models := []model.CatalogModel{{Name: "model-a"}}
	keys := collectCustomPropertyKeys(models)
	assert.Empty(t, keys)
}

func TestBuildCSVHeaders(t *testing.T) {
	headers := buildCSVHeaders([]string{"arch", "version"})
	expected := append(fixedCSVColumns, "custom_arch", "custom_version")
	assert.Equal(t, expected, headers)
}

func TestBuildCSVHeadersNoCustom(t *testing.T) {
	headers := buildCSVHeaders(nil)
	assert.Equal(t, fixedCSVColumns, headers)
}

func TestMetadataValueToString(t *testing.T) {
	tests := []struct {
		name     string
		value    model.MetadataValue
		expected string
	}{
		{
			name:     "string value",
			value:    model.MetadataStringValueAsMetadataValue(model.NewMetadataStringValue("hello", "MetadataStringValue")),
			expected: "hello",
		},
		{
			name:     "int value",
			value:    model.MetadataIntValueAsMetadataValue(model.NewMetadataIntValue("42", "MetadataIntValue")),
			expected: "42",
		},
		{
			name:     "double value",
			value:    model.MetadataDoubleValueAsMetadataValue(model.NewMetadataDoubleValue(3.14, "MetadataDoubleValue")),
			expected: "3.14",
		},
		{
			name:     "bool value",
			value:    model.MetadataBoolValueAsMetadataValue(model.NewMetadataBoolValue(true, "MetadataBoolValue")),
			expected: "true",
		},
		{
			name:     "struct value",
			value:    model.MetadataStructValueAsMetadataValue(model.NewMetadataStructValue(`{"key":"val"}`, "MetadataStructValue")),
			expected: `{"key":"val"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, metadataValueToString(tc.value))
		})
	}
}

func TestModelToCSVRow(t *testing.T) {
	m := model.CatalogModel{
		Id:                      strPtr("123"),
		Name:                    "granite-8b",
		SourceId:                strPtr("huggingface"),
		Description:             strPtr("A model"),
		Provider:                strPtr("IBM"),
		License:                 strPtr("Apache-2.0"),
		LicenseLink:             strPtr("https://example.com"),
		LibraryName:             strPtr("transformers"),
		Maturity:                strPtr("GA"),
		Language:                []string{"en", "fr"},
		Tasks:                   []string{"text-generation", "summarization"},
		CreateTimeSinceEpoch:    strPtr("1700000000000"),
		LastUpdateTimeSinceEpoch: strPtr("1700100000000"),
		CustomProperties: map[string]model.MetadataValue{
			"arch": model.MetadataStringValueAsMetadataValue(model.NewMetadataStringValue("transformer", "MetadataStringValue")),
		},
	}

	customKeys := []string{"arch", "version"}
	row := modelToCSVRow(m, customKeys)

	assert.Equal(t, "123", row[0])
	assert.Equal(t, "granite-8b", row[1])
	assert.Equal(t, "huggingface", row[2])
	assert.Equal(t, "A model", row[3])
	assert.Equal(t, "IBM", row[4])
	assert.Equal(t, "Apache-2.0", row[5])
	assert.Equal(t, "https://example.com", row[6])
	assert.Equal(t, "transformers", row[7])
	assert.Equal(t, "GA", row[8])
	assert.Equal(t, "en;fr", row[9])
	assert.Equal(t, "text-generation;summarization", row[10])
	assert.Equal(t, "1700000000000", row[11])
	assert.Equal(t, "1700100000000", row[12])
	assert.Equal(t, "transformer", row[13]) // custom_arch
	assert.Equal(t, "", row[14])            // custom_version (missing)
}

func TestModelToCSVRowNilFields(t *testing.T) {
	m := model.CatalogModel{Name: "minimal-model"}
	row := modelToCSVRow(m, nil)

	assert.Equal(t, "", row[0])              // id
	assert.Equal(t, "minimal-model", row[1]) // name
	assert.Equal(t, "", row[2])              // source_id
}

func TestWriteCSV(t *testing.T) {
	models := []model.CatalogModel{
		{
			Id:       strPtr("1"),
			Name:     "model-a",
			Provider: strPtr("IBM"),
			Language: []string{"en"},
			CustomProperties: map[string]model.MetadataValue{
				"arch": model.MetadataStringValueAsMetadataValue(model.NewMetadataStringValue("transformer", "MetadataStringValue")),
			},
		},
		{
			Id:       strPtr("2"),
			Name:     "model-b",
			Provider: strPtr("Meta"),
		},
	}

	var buf bytes.Buffer
	err := writeCSV(&buf, models)
	require.NoError(t, err)

	content := buf.String()

	// Check BOM
	assert.True(t, strings.HasPrefix(content, "\xEF\xBB\xBF"), "should start with UTF-8 BOM")

	// Parse CSV (skip BOM)
	reader := csv.NewReader(strings.NewReader(content[3:]))
	records, err := reader.ReadAll()
	require.NoError(t, err)

	// Header + 2 data rows
	require.Len(t, records, 3)

	// Check header includes custom column
	header := records[0]
	assert.Contains(t, header, "custom_arch")

	// Check data rows
	assert.Equal(t, "1", records[1][0])
	assert.Equal(t, "model-a", records[1][1])
	assert.Equal(t, "2", records[2][0])
	assert.Equal(t, "model-b", records[2][1])
}

func TestWriteCSVEmpty(t *testing.T) {
	var buf bytes.Buffer
	err := writeCSV(&buf, []model.CatalogModel{})
	require.NoError(t, err)

	content := buf.String()
	assert.True(t, strings.HasPrefix(content, "\xEF\xBB\xBF"))

	reader := csv.NewReader(strings.NewReader(content[3:]))
	records, err := reader.ReadAll()
	require.NoError(t, err)

	// Header only, no data rows
	require.Len(t, records, 1)
	assert.Equal(t, fixedCSVColumns, records[0])
}

func TestWriteCSVToFile(t *testing.T) {
	models := []model.CatalogModel{
		{Id: strPtr("1"), Name: "model-a"},
	}

	f, err := writeCSVToFile(models)
	require.NoError(t, err)
	defer f.Close()

	assert.True(t, strings.HasSuffix(f.Name(), ".csv"))
	defer func() { _ = os.Remove(f.Name()) }()

	var buf bytes.Buffer
	_, err = buf.ReadFrom(f)
	require.NoError(t, err)

	assert.True(t, strings.HasPrefix(buf.String(), "\xEF\xBB\xBF"))
}
