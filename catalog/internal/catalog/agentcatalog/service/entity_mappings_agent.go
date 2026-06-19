package service

import "github.com/kubeflow/hub/internal/platform/db/filter"

type agentEntityMappings struct{}

// Unexported: only used by the repository constructor in this package.
// Export (New...) if tests in a parent package need to call it directly.
func newAgentEntityMappings() filter.EntityMappingFunctions {
	return &agentEntityMappings{}
}

func (m *agentEntityMappings) GetMLMDEntityType(_ filter.RestEntityType) filter.EntityType {
	return filter.EntityTypeContext
}

func (m *agentEntityMappings) GetPropertyDefinitionForRestEntity(_ filter.RestEntityType, propertyName string) filter.PropertyDefinition {
	if def, ok := agentProperties[propertyName]; ok {
		return def
	}
	return filter.PropertyDefinition{
		Location:  filter.Custom,
		ValueType: filter.StringValueType,
		Column:    propertyName,
	}
}

func (m *agentEntityMappings) IsChildEntity(_ filter.RestEntityType) bool {
	return false
}

var agentProperties = map[string]filter.PropertyDefinition{
	"id":                       {Location: filter.EntityTable, ValueType: filter.IntValueType, Column: "id"},
	"name":                     {Location: filter.EntityTable, ValueType: filter.StringValueType, Column: "name"},
	"externalId":               {Location: filter.EntityTable, ValueType: filter.StringValueType, Column: "external_id"},
	"createTimeSinceEpoch":     {Location: filter.EntityTable, ValueType: filter.IntValueType, Column: "create_time_since_epoch"},
	"lastUpdateTimeSinceEpoch": {Location: filter.EntityTable, ValueType: filter.IntValueType, Column: "last_update_time_since_epoch"},
	"source_id":                {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "source_id"},
	"description":              {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "description"},
	"displayName":              {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "displayName"},
	"framework":                {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "framework"},
	"agentType":                {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "agentType"},
	"tags":                     {Location: filter.PropertyTable, ValueType: filter.ArrayValueType, Column: "tags"},
	"models":                   {Location: filter.PropertyTable, ValueType: filter.ArrayValueType, Column: "models"},
	"logo":                     {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "logo"},
	"repositoryUrl":            {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "repositoryUrl"},
	"publishedDate":            {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "publishedDate"},
	"readme":                   {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "readme"},
	"env":                      {Location: filter.PropertyTable, ValueType: filter.ArrayValueType, Column: "env"},
	"artifacts":                {Location: filter.PropertyTable, ValueType: filter.ArrayValueType, Column: "artifacts"},
}
