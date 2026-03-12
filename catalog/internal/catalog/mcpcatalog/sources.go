package mcpcatalog

import (
	"github.com/kubeflow/model-registry/catalog/internal/catalog/basecatalog"
	model "github.com/kubeflow/model-registry/catalog/pkg/openapi"
	"github.com/kubeflow/model-registry/internal/apiutils"
)

// MCPSourceCollection manages MCP catalog sources from multiple origins with priority-based merging.
// It wraps the generic GenericSourceCollection and adds MCP-specific convenience methods.
type MCPSourceCollection struct {
	*basecatalog.GenericSourceCollection[basecatalog.MCPSource]
}

// NewMCPSourceCollection creates a new MCPSourceCollection with the given origin order.
// Origins listed later in the order take precedence over earlier ones.
func NewMCPSourceCollection(originOrder ...string) *MCPSourceCollection {
	return &MCPSourceCollection{
		GenericSourceCollection: basecatalog.NewGenericSourceCollection(
			mergeMCPSources,
			applyMCPDefaults,
			originOrder...,
		),
	}
}

// mergeMCPSources performs field-level merging of two MCPSource structs.
// Fields from 'override' take precedence over 'base' when they are explicitly set.
// A field is considered "set" if:
// - For strings: non-empty
// - For pointers: non-nil
// - For slices: non-nil (empty slice is considered explicitly set to "no items")
// - For maps: non-nil (empty map is considered explicitly set)
func mergeMCPSources(base, override basecatalog.MCPSource) basecatalog.MCPSource {
	result := base

	// ID is always taken from override (it's the key)
	result.ID = override.ID

	// Merge shared fields using the common helper
	common := basecatalog.MergeCommonSourceFields(
		basecatalog.CommonSourceFields{Name: base.Name, Enabled: base.Enabled, Labels: base.Labels, Type: base.Type, Properties: base.Properties, Origin: base.Origin, AssetType: base.AssetType},
		basecatalog.CommonSourceFields{Name: override.Name, Enabled: override.Enabled, Labels: override.Labels, Type: override.Type, Properties: override.Properties, Origin: override.Origin, AssetType: override.AssetType},
	)
	result.Name = common.Name
	result.Enabled = common.Enabled
	result.Labels = common.Labels
	result.Type = common.Type
	result.Properties = common.Properties
	result.Origin = common.Origin
	result.AssetType = common.AssetType

	return result
}

// applyMCPDefaults applies default values to an MCPSource for fields that are not set.
func applyMCPDefaults(source basecatalog.MCPSource) basecatalog.MCPSource {
	if source.Enabled == nil {
		source.Enabled = apiutils.Of(true)
	}
	if source.Labels == nil {
		source.Labels = []string{}
	}
	if source.AssetType == nil {
		source.AssetType = model.CATALOGASSETTYPE_MCP_SERVERS.Ptr()
	}
	return source
}
