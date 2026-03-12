package filter

import (
	"strings"

	"github.com/kubeflow/model-registry/catalog/internal/catalog/basecatalog"
	catalogmodels "github.com/kubeflow/model-registry/catalog/internal/db/models"
	"github.com/kubeflow/model-registry/internal/db/filter"
)

// catalogEntityMappings wraps a CatalogEntityRegistry and adds catalog-specific
// behaviors like equality expansion for externalId.
type catalogEntityMappings struct {
	*basecatalog.CatalogEntityRegistry
}

// NewCatalogEntityMappings creates a new instance of catalog entity mappings
// using the declarative CatalogEntityRegistry.
func NewCatalogEntityMappings() filter.EntityMappingFunctions {
	return &catalogEntityMappings{
		CatalogEntityRegistry: buildCatalogRegistry(),
	}
}

// GetEqualityExpansion implements filter.EqualityExpander so externalId = "x" matches both
// exact and namespaced (sourceId:x) stored values, returning all models regardless of source.
func (c *catalogEntityMappings) GetEqualityExpansion(restEntityType filter.RestEntityType, propertyName string, value any) (likeArg any, useExpansion bool) {
	if restEntityType != filter.RestEntityType(catalogmodels.RestEntityCatalogModel) || propertyName != "externalId" {
		return nil, false
	}
	strVal, ok := value.(string)
	if !ok || strVal == "" {
		return nil, false
	}
	// Match any source prefix: sourceId:externalId
	return "%:" + escapeLike(strVal), true
}

// escapeLike escapes SQL LIKE metacharacters (%, _, \) for safe use as a literal in a LIKE pattern.
func escapeLike(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch r {
		case '\\', '%', '_':
			b.WriteRune('\\')
			b.WriteRune(r)
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

// buildCatalogRegistry creates and populates a CatalogEntityRegistry with all
// catalog entity type definitions. This replaces the previous hand-coded
// switch/if chains with declarative registration.
func buildCatalogRegistry() *basecatalog.CatalogEntityRegistry {
	reg := basecatalog.NewCatalogEntityRegistry()

	// CatalogModel (Context entity)
	reg.Register(filter.RestEntityType(catalogmodels.RestEntityCatalogModel), basecatalog.EntityTypeDefinition{
		MLMDEntityType: filter.EntityTypeContext,
		Properties: basecatalog.MergeProperties(
			basecatalog.CommonContextProperties(),
			basecatalog.CommonCatalogContextProperties(),
			map[string]filter.PropertyDefinition{
				"owner":          {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "owner"},
				"state":          {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "state"},
				"language":       {Location: filter.PropertyTable, ValueType: filter.ArrayValueType, Column: "language"},
				"library_name":   {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "library_name"},
				"maturity":       {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "maturity"},
				"tasks":          {Location: filter.PropertyTable, ValueType: filter.ArrayValueType, Column: "tasks"},
				"tags":           {Location: filter.PropertyTable, ValueType: filter.ArrayValueType, Column: "tags"},
				"verifiedSource": {Location: filter.PropertyTable, ValueType: filter.BoolValueType, Column: "verifiedSource"},
			},
		),
		IsChild:                false,
		RelatedEntityPrefix:    "artifacts.",
		RelatedEntityType:      filter.RelatedEntityArtifact,
		RelatedEntityJoinTable: "Attribution",
	})

	// CatalogArtifact (Artifact entity)
	reg.Register(filter.RestEntityType(catalogmodels.RestEntityCatalogArtifact), basecatalog.EntityTypeDefinition{
		MLMDEntityType: filter.EntityTypeArtifact,
		Properties: basecatalog.MergeProperties(
			basecatalog.CommonArtifactProperties(),
			map[string]filter.PropertyDefinition{
				"artifactType": {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "artifactType"},
			},
		),
		IsChild: false,
	})

	// MCPServer (Context entity)
	reg.Register(filter.RestEntityType(catalogmodels.RestEntityMCPServer), basecatalog.EntityTypeDefinition{
		MLMDEntityType: filter.EntityTypeContext,
		Properties: basecatalog.MergeProperties(
			basecatalog.CommonContextProperties(),
			basecatalog.CommonCatalogContextProperties(),
			map[string]filter.PropertyDefinition{
				"base_name":        {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "base_name"},
				"name":             {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "base_name"},
				"version":          {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "version"},
				"tags":             {Location: filter.PropertyTable, ValueType: filter.ArrayValueType, Column: "tags"},
				"transports":       {Location: filter.PropertyTable, ValueType: filter.ArrayValueType, Column: "transports"},
				"deploymentMode":   {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "deploymentMode"},
				"documentationUrl": {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "documentationUrl"},
				"repositoryUrl":    {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "repositoryUrl"},
				"sourceCode":       {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "sourceCode"},
				"publishedDate":    {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "publishedDate"},
				"lastUpdated":      {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "lastUpdated"},
				"verifiedSource":   {Location: filter.PropertyTable, ValueType: filter.BoolValueType, Column: "verifiedSource"},
				"secureEndpoint":   {Location: filter.PropertyTable, ValueType: filter.BoolValueType, Column: "secureEndpoint"},
				"sast":             {Location: filter.PropertyTable, ValueType: filter.BoolValueType, Column: "sast"},
				"readOnlyTools":    {Location: filter.PropertyTable, ValueType: filter.BoolValueType, Column: "readOnlyTools"},
				"endpoints":        {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "endpoints"},
				"artifacts":        {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "artifacts"},
				"runtimeMetadata":  {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "runtimeMetadata"},
			},
		),
		IsChild: false,
	})

	// MCPServerTool (Execution entity)
	reg.Register(filter.RestEntityType(catalogmodels.RestEntityMCPServerTool), basecatalog.EntityTypeDefinition{
		MLMDEntityType: filter.EntityTypeExecution,
		Properties: basecatalog.MergeProperties(
			map[string]filter.PropertyDefinition{
				"id":                       {Location: filter.EntityTable, ValueType: filter.IntValueType, Column: "id"},
				"name":                     {Location: filter.EntityTable, ValueType: filter.StringValueType, Column: "name"},
				"createTimeSinceEpoch":     {Location: filter.EntityTable, ValueType: filter.IntValueType, Column: "create_time_since_epoch"},
				"lastUpdateTimeSinceEpoch": {Location: filter.EntityTable, ValueType: filter.IntValueType, Column: "last_update_time_since_epoch"},
			},
			map[string]filter.PropertyDefinition{
				"description": {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "description"},
				"accessType":  {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "accessType"},
			},
		),
		IsChild: false,
	})

	return reg
}
