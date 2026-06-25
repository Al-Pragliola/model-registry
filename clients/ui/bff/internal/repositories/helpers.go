package repositories

import (
	"fmt"
	"net/url"

	"github.com/kubeflow/hub/ui/bff/internal/models"
	"gopkg.in/yaml.v3"
)

func FilterPageValues(values url.Values) url.Values {
	result := url.Values{}

	if v := values.Get("pageSize"); v != "" {
		result.Set("pageSize", v)
	}
	if v := values.Get("orderBy"); v != "" {
		result.Set("orderBy", v)
	}
	if v := values.Get("sortOrder"); v != "" {
		result.Set("sortOrder", v)
	}
	if v := values.Get("nextPageToken"); v != "" {
		result.Set("nextPageToken", v)
	}
	if v := values.Get("name"); v != "" {
		result.Set("name", v)
	}
	if v := values.Get("q"); v != "" {
		result.Set("q", v)
	}
	if v := values.Get("source"); v != "" {
		result.Set("source", v)
	}
	if v := values.Get("sourceLabel"); v != "" {
		result.Set("sourceLabel", v)
	}
	if v := values.Get("filterQuery"); v != "" {
		result.Set("filterQuery", v)
	}
	if v := values.Get("artifactType"); v != "" {
		result.Set("artifactType", v)
	}
	if v := values.Get("targetRPS"); v != "" {
		result.Set("targetRPS", v)
	}
	if v := values.Get("recommendations"); v != "" {
		result.Set("recommendations", v)
	}
	if v := values.Get("rpsProperty"); v != "" {
		result.Set("rpsProperty", v)
	}
	if v := values.Get("latencyProperty"); v != "" {
		result.Set("latencyProperty", v)
	}
	if v := values.Get("hardwareCountProperty"); v != "" {
		result.Set("hardwareCountProperty", v)
	}
	if v := values.Get("hardwareTypeProperty"); v != "" {
		result.Set("hardwareTypeProperty", v)
	}
	if v := values.Get("filterStatus"); v != "" {
		result.Set("filterStatus", v)
	}
	if v := values.Get("assetType"); v != "" {
		result.Set("assetType", v)
	}
	if v := values.Get("includeTools"); v != "" {
		result.Set("includeTools", v)
	}

	return result
}

func UrlWithParams(url string, values url.Values) string {
	queryString := values.Encode()
	if queryString == "" {
		return url
	}
	return fmt.Sprintf("%s?%s", url, queryString)
}

func UrlWithPageParams(url string, values url.Values) string {
	pageValues := FilterPageValues(values)
	return UrlWithParams(url, pageValues)
}

const (
	SectionKeyCatalogs      = "catalogs"
	SectionKeyAgentCatalogs = "agent_catalogs"
)

func ParseCatalogYaml(raw string, isDefault bool) ([]models.CatalogSourceConfig, error) {
	return ParseCatalogYamlSection(raw, isDefault, SectionKeyCatalogs)
}

func ParseCatalogYamlSection(raw string, isDefault bool, sectionKey string) ([]models.CatalogSourceConfig, error) {
	var fullDoc map[string]interface{}
	if err := yaml.Unmarshal([]byte(raw), &fullDoc); err != nil {
		return nil, fmt.Errorf("failed to parse yaml: %w", err)
	}

	sectionRaw, ok := fullDoc[sectionKey]
	if !ok {
		return []models.CatalogSourceConfig{}, nil
	}

	sectionList, ok := sectionRaw.([]interface{})
	if !ok {
		return []models.CatalogSourceConfig{}, nil
	}

	catalogs := make([]models.CatalogSourceConfig, 0, len(sectionList))
	for _, item := range sectionList {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		entry := models.CatalogSourceConfig{
			Id:        getStringFromMap(itemMap, "id"),
			Name:      getStringFromMap(itemMap, "name"),
			Type:      getStringFromMap(itemMap, "type"),
			IsDefault: &isDefault,
		}

		if enabled, ok := itemMap["enabled"]; ok {
			if b, ok := enabled.(bool); ok {
				entry.Enabled = &b
			}
		}

		if labels, ok := itemMap["labels"]; ok {
			entry.Labels = ExtractStringSlice(labels)
		}

		entry.IncludedModels = ExtractStringSlice(itemMap["includedModels"])
		entry.ExcludedModels = ExtractStringSlice(itemMap["excludedModels"])
		entry.IncludedAgents = ExtractStringSlice(itemMap["includedAgents"])
		entry.ExcludedAgents = ExtractStringSlice(itemMap["excludedAgents"])

		if props, ok := itemMap["properties"].(map[string]interface{}); ok {
			if allowedOrganization, ok := props["allowedOrganization"].(string); ok {
				entry.AllowedOrganization = &allowedOrganization
			}
		}

		catalogs = append(catalogs, entry)
	}

	return catalogs, nil
}

func getStringFromMap(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func ExtractStringSlice(value interface{}) []string {
	if arr, ok := value.([]interface{}); ok {
		result := make([]string, 0, len(arr))
		for _, item := range arr {
			if str, ok := item.(string); ok {
				result = append(result, str)
			}
		}
		return result
	}
	if strSlice, ok := value.([]string); ok {
		return strSlice
	}
	return []string{}

}

func FindCatalogSourceById(sourceYAML string, catalogId string, isDefault bool) *models.CatalogSourceConfig {
	return FindCatalogSourceByIdSection(sourceYAML, catalogId, isDefault, SectionKeyCatalogs)
}

func FindCatalogSourceByIdSection(sourceYAML string, catalogId string, isDefault bool, sectionKey string) *models.CatalogSourceConfig {
	if sourceYAML == "" {
		return nil
	}

	configs, err := ParseCatalogYamlSection(sourceYAML, isDefault, sectionKey)
	if err != nil {
		return nil
	}

	for _, config := range configs {
		if config.Id == catalogId {
			c := config
			return &c
		}
	}

	return nil
}

func ConvertSourceConfigToYamlEntry(payload models.CatalogSourceConfigPayload,
	yamlFileName string,
	secretName string) map[string]interface{} {
	entry := map[string]interface{}{
		"id":      payload.Id,
		"name":    payload.Name,
		"type":    payload.Type,
		"enabled": payload.Enabled,
	}

	if len(payload.Labels) > 0 {
		entry["labels"] = payload.Labels
	}

	properties := make(map[string]interface{})

	switch payload.Type {
	case CatalogTypeYaml:
		properties["yamlCatalogPath"] = yamlFileName
	case CatalogTypeHuggingFace:
		properties[ApiKey] = secretName
		if payload.AllowedOrganization != nil {
			properties["allowedOrganization"] = *payload.AllowedOrganization
		}

	}

	if len(properties) > 0 {
		entry["properties"] = properties
	}

	if len(payload.IncludedModels) > 0 {
		entry["includedModels"] = payload.IncludedModels
	}
	if len(payload.ExcludedModels) > 0 {
		entry["excludedModels"] = payload.ExcludedModels
	}
	if len(payload.IncludedAgents) > 0 {
		entry["includedAgents"] = payload.IncludedAgents
	}
	if len(payload.ExcludedAgents) > 0 {
		entry["excludedAgents"] = payload.ExcludedAgents
	}

	return entry
}

func AppendCatalogSourceToYaml(existingConfigMapEntry string, newEntry map[string]interface{}) (string, error) {
	return AppendCatalogSourceToYamlSection(existingConfigMapEntry, newEntry, SectionKeyCatalogs)
}

func AppendCatalogSourceToYamlSection(existingConfigMapEntry string, newEntry map[string]interface{}, sectionKey string) (string, error) {
	var fullDoc map[string]interface{}
	if existingConfigMapEntry != "" {
		if err := yaml.Unmarshal([]byte(existingConfigMapEntry), &fullDoc); err != nil {
			return "", fmt.Errorf("failed to parse existing sources.yaml: %w", err)
		}
	}
	if fullDoc == nil {
		fullDoc = make(map[string]interface{})
	}

	var sectionList []interface{}
	if existing, ok := fullDoc[sectionKey]; ok {
		if list, ok := existing.([]interface{}); ok {
			sectionList = list
		}
	}
	sectionList = append(sectionList, newEntry)
	fullDoc[sectionKey] = sectionList

	updatedBytes, err := yaml.Marshal(fullDoc)
	if err != nil {
		return "", fmt.Errorf("failed to marshal updated sources.yaml: %w", err)
	}

	return string(updatedBytes), nil
}

func RemoveCatalogSourceFromYAML(existingYAML string, sourceId string) (string, error) {
	return RemoveCatalogSourceFromYAMLSection(existingYAML, sourceId, SectionKeyCatalogs)
}

func RemoveCatalogSourceFromYAMLSection(existingYAML string, sourceId string, sectionKey string) (string, error) {
	var fullDoc map[string]interface{}
	if err := yaml.Unmarshal([]byte(existingYAML), &fullDoc); err != nil {
		return "", fmt.Errorf("failed to parse sources.yaml: %w", err)
	}

	sectionRaw, ok := fullDoc[sectionKey]
	if !ok {
		return existingYAML, nil
	}
	sectionList, ok := sectionRaw.([]interface{})
	if !ok {
		return existingYAML, nil
	}

	filtered := make([]interface{}, 0)
	for _, item := range sectionList {
		if itemMap, ok := item.(map[string]interface{}); ok {
			if id, ok := itemMap["id"].(string); ok && id != sourceId {
				filtered = append(filtered, item)
			}
		}
	}

	fullDoc[sectionKey] = filtered
	updatedBytes, err := yaml.Marshal(fullDoc)
	if err != nil {
		return "", fmt.Errorf("failed to marshal updated sources.yaml: %w", err)
	}

	return string(updatedBytes), nil
}

func FindCatalogSourceProperties(sourceYAML string, sourceId string) (secretName string, yamlPath string) {
	return FindCatalogSourcePropertiesSection(sourceYAML, sourceId, SectionKeyCatalogs)
}

func FindCatalogSourcePropertiesSection(sourceYAML string, sourceId string, sectionKey string) (secretName string, yamlPath string) {
	if sourceYAML == "" {
		return "", ""
	}

	var fullDoc map[string]interface{}
	if err := yaml.Unmarshal([]byte(sourceYAML), &fullDoc); err != nil {
		return "", ""
	}

	sectionRaw, ok := fullDoc[sectionKey]
	if !ok {
		return "", ""
	}
	sectionList, ok := sectionRaw.([]interface{})
	if !ok {
		return "", ""
	}

	for _, item := range sectionList {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		id, _ := itemMap["id"].(string)
		if id == sourceId {
			if props, ok := itemMap["properties"].(map[string]interface{}); ok {
				secretName, _ = props[ApiKey].(string)
				yamlPath, _ = props["yamlCatalogPath"].(string)
			}
			return
		}
	}
	return "", ""
}

func UpdateCatalogSourceInYAML(
	existingYAML string,
	catalogId string,
	payload models.CatalogSourceConfigPayload,
	secretName string,
	yamlFilePath string,
) (string, error) {
	return UpdateCatalogSourceInYAMLSection(existingYAML, catalogId, payload, secretName, yamlFilePath, SectionKeyCatalogs)
}

func UpdateCatalogSourceInYAMLSection(
	existingYAML string,
	catalogId string,
	payload models.CatalogSourceConfigPayload,
	secretName string,
	yamlFilePath string,
	sectionKey string,
) (string, error) {
	if existingYAML == "" {
		return "", fmt.Errorf("no existing yaml to update")
	}

	var fullDoc map[string]interface{}
	if err := yaml.Unmarshal([]byte(existingYAML), &fullDoc); err != nil {
		return "", fmt.Errorf("failed to parse sources.yaml: %w", err)
	}

	sectionRaw, ok := fullDoc[sectionKey]
	if !ok {
		return "", fmt.Errorf("section '%s' not found in yaml", sectionKey)
	}
	sectionList, ok := sectionRaw.([]interface{})
	if !ok {
		return "", fmt.Errorf("section '%s' is not a list", sectionKey)
	}

	found := false
	for i, item := range sectionList {
		catalogSource, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		if id, ok := catalogSource["id"].(string); ok && id == catalogId {
			found = true

			if payload.Name != "" {
				catalogSource["name"] = payload.Name
			}
			if len(payload.Labels) > 0 {
				catalogSource["labels"] = payload.Labels
			}
			if payload.Enabled != nil {
				catalogSource["enabled"] = *payload.Enabled
			}

			properties, _ := catalogSource["properties"].(map[string]interface{})
			if properties == nil {
				properties = make(map[string]interface{})
			}

			if payload.IncludedModels != nil {
				if len(payload.IncludedModels) > 0 {
					catalogSource["includedModels"] = payload.IncludedModels
				} else {
					delete(catalogSource, "includedModels")
				}
			}
			if payload.ExcludedModels != nil {
				if len(payload.ExcludedModels) > 0 {
					catalogSource["excludedModels"] = payload.ExcludedModels
				} else {
					delete(catalogSource, "excludedModels")
				}
			}
			if payload.IncludedAgents != nil {
				if len(payload.IncludedAgents) > 0 {
					catalogSource["includedAgents"] = payload.IncludedAgents
				} else {
					delete(catalogSource, "includedAgents")
				}
			}
			if payload.ExcludedAgents != nil {
				if len(payload.ExcludedAgents) > 0 {
					catalogSource["excludedAgents"] = payload.ExcludedAgents
				} else {
					delete(catalogSource, "excludedAgents")
				}
			}
			if payload.AllowedOrganization != nil {
				properties["allowedOrganization"] = *payload.AllowedOrganization
			}
			if secretName != "" {
				properties["apiKey"] = secretName
			}
			if yamlFilePath != "" && payload.Yaml != nil {
				properties["yamlCatalogPath"] = yamlFilePath
			}

			if len(properties) > 0 {
				catalogSource["properties"] = properties
			}

			sectionList[i] = catalogSource
			break
		}
	}

	if !found {
		return "", fmt.Errorf("catalog '%s' not found in yaml", catalogId)
	}

	fullDoc[sectionKey] = sectionList
	updatedBytes, err := yaml.Marshal(fullDoc)
	if err != nil {
		return "", fmt.Errorf("failed to marshal updated sources.yaml: %w", err)
	}

	return string(updatedBytes), nil
}

func BuildOverrideEntryForDefaultSource(catalogId string, payload models.CatalogSourceConfigPayload) map[string]interface{} {
	entry := map[string]interface{}{
		"id": catalogId,
	}

	if payload.Enabled != nil {
		entry["enabled"] = *payload.Enabled
	}

	if len(payload.IncludedModels) > 0 {
		entry["includedModels"] = payload.IncludedModels
	}
	if len(payload.ExcludedModels) > 0 {
		entry["excludedModels"] = payload.ExcludedModels
	}
	if len(payload.IncludedAgents) > 0 {
		entry["includedAgents"] = payload.IncludedAgents
	}
	if len(payload.ExcludedAgents) > 0 {
		entry["excludedAgents"] = payload.ExcludedAgents
	}

	return entry
}
