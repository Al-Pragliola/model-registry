package modelcatalog

import (
	"maps"
	"slices"
	"strings"

	"github.com/kubeflow/model-registry/catalog/internal/catalog/basecatalog"
	model "github.com/kubeflow/model-registry/catalog/pkg/openapi"
	"github.com/kubeflow/model-registry/internal/apiutils"
)

// SourceCollection manages model catalog sources from multiple origins with priority-based merging.
// It wraps the generic GenericSourceCollection and adds model-specific convenience methods.
type SourceCollection struct {
	*basecatalog.GenericSourceCollection[basecatalog.ModelSource]
}

// NewSourceCollection creates a new SourceCollection with the given origin order.
// Origins listed later in the order take precedence over earlier ones.
func NewSourceCollection(originOrder ...string) *SourceCollection {
	return &SourceCollection{
		GenericSourceCollection: basecatalog.NewGenericSourceCollection(
			mergeSources,
			applyDefaults,
			originOrder...,
		),
	}
}

// mergeSources performs field-level merging of two Source structs.
// Fields from 'override' take precedence over 'base' when they are explicitly set.
// A field is considered "set" if:
// - For strings: non-empty
// - For pointers: non-nil
// - For slices: non-nil (empty slice is considered explicitly set to "no items")
// - For maps: non-nil (empty map is considered explicitly set)
func mergeSources(base, override basecatalog.ModelSource) basecatalog.ModelSource {
	result := base

	// Id is always taken from override (it's the key)
	result.Id = override.Id

	// Merge shared fields using the common helper
	common := basecatalog.MergeCommonSourceFields(
		basecatalog.CommonSourceFields{Name: base.Name, Enabled: base.Enabled, Labels: base.Labels, Type: base.Type, Properties: base.Properties, Origin: base.Origin},
		basecatalog.CommonSourceFields{Name: override.Name, Enabled: override.Enabled, Labels: override.Labels, Type: override.Type, Properties: override.Properties, Origin: override.Origin},
	)
	result.Name = common.Name
	result.Enabled = common.Enabled
	result.Labels = common.Labels
	result.Type = common.Type
	result.Properties = common.Properties
	result.Origin = common.Origin

	// Model-specific fields
	if override.IncludedModels != nil {
		result.IncludedModels = override.IncludedModels
	}
	if override.ExcludedModels != nil {
		result.ExcludedModels = override.ExcludedModels
	}

	return result
}

// applyDefaults applies default values to an Source for fields that are not set.
func applyDefaults(source basecatalog.ModelSource) basecatalog.ModelSource {
	// Default Enabled to true if not set
	if source.Enabled == nil {
		source.Enabled = apiutils.Of(true)
	}

	// Default Labels to empty slice if not set
	if source.Labels == nil {
		source.Labels = []string{}
	}

	return source
}

// All returns all sources as CatalogSource (for the API).
// This excludes internal fields like Type and Properties.
func (sc *SourceCollection) All() map[string]model.CatalogSource {
	result := map[string]model.CatalogSource{}
	for id, source := range sc.AllSources() {
		result[id] = source.CatalogSource
	}
	return result
}

// Get returns a source by name if it exists and is enabled.
func (sc *SourceCollection) Get(name string) (src model.CatalogSource, ok bool) {
	allSources := sc.AllSources()

	source, exists := allSources[name]
	if !exists {
		return model.CatalogSource{}, false
	}

	// Only return if enabled
	if source.Enabled != nil && *source.Enabled {
		return source.CatalogSource, true
	}
	return model.CatalogSource{}, false
}

// ByLabel returns enabled sources that have any of the labels provided. The matching
// is case insensitive.
//
// If a label is "null", every source without a label is returned.
func (sc *SourceCollection) ByLabel(labels []string) []model.CatalogSource {
	labelMap := make(map[string]struct{}, len(labels))
	for _, label := range labels {
		labelMap[strings.ToLower(label)] = struct{}{}
	}

	matches := map[string]model.CatalogSource{}
	sources := sc.AllSources()

	if _, hasNull := labelMap["null"]; hasNull {
		for _, source := range sources {
			// Skip disabled sources
			if source.Enabled == nil || !*source.Enabled {
				continue
			}
			if len(source.Labels) == 0 {
				matches[source.Id] = source.CatalogSource
			}
		}
	}

OUTER:
	for _, source := range sources {
		// Skip disabled sources
		if source.Enabled == nil || !*source.Enabled {
			continue
		}
		for _, label := range source.Labels {
			if _, match := labelMap[strings.ToLower(label)]; match {
				matches[source.Id] = source.CatalogSource
				continue OUTER
			}
		}
	}

	return slices.Collect(maps.Values(matches))
}
