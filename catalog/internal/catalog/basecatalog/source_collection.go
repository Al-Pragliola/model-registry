package basecatalog

import (
	"maps"
	"sync"
)

// Source is the constraint for types usable in a GenericSourceCollection.
// Any catalog source struct that has a string ID satisfies this.
type Source interface {
	GetId() string
}

// SourceMergeFunc performs field-level merging of two sources of the same type.
// Fields from override take precedence over base.
type SourceMergeFunc[S any] func(base, override S) S

// SourceDefaultsFunc applies default values to a source for unset fields.
type SourceDefaultsFunc[S any] func(S) S

// originEntry holds sources from a single origin (config file).
type originEntry[S any] struct {
	origin  string
	sources map[string]S
}

// GenericSourceCollection manages catalog sources from multiple origins with priority-based merging.
// Later entries in the entries slice take precedence over earlier ones.
//
// This is the generic implementation used by both ModelSource and MCPSource collections.
// Domain-specific convenience methods (e.g. ByLabel, Get) should be added by wrapper types.
type GenericSourceCollection[S Source] struct {
	mu           sync.RWMutex
	entries      []originEntry[S]
	namedQueries map[string]map[string]FieldFilter

	// Pluggable behaviors
	mergeFunc    SourceMergeFunc[S]
	defaultsFunc SourceDefaultsFunc[S]
}

// NewGenericSourceCollection creates a new GenericSourceCollection with the given origin order
// and pluggable merge/defaults functions.
// Origins listed later in the order take precedence over earlier ones.
func NewGenericSourceCollection[S Source](
	mergeFunc SourceMergeFunc[S],
	defaultsFunc SourceDefaultsFunc[S],
	originOrder ...string,
) *GenericSourceCollection[S] {
	entries := make([]originEntry[S], len(originOrder))
	for i, origin := range originOrder {
		entries[i] = originEntry[S]{origin: origin, sources: nil}
	}
	return &GenericSourceCollection[S]{
		entries:      entries,
		namedQueries: make(map[string]map[string]FieldFilter),
		mergeFunc:    mergeFunc,
		defaultsFunc: defaultsFunc,
	}
}

// Merge adds sources from one origin, completely replacing anything that was
// previously from that origin.
//
// If a source with the same ID exists in multiple origins, fields from
// higher-priority origins (listed later in entries) override fields from
// lower-priority origins via the configured mergeFunc.
func (sc *GenericSourceCollection[S]) Merge(origin string, sources map[string]S) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	return sc.mergeSourcesInternal(origin, sources)
}

// MergeWithNamedQueries adds sources and named queries from one origin.
// Later origins override earlier ones at the field level within a query.
func (sc *GenericSourceCollection[S]) MergeWithNamedQueries(origin string, sources map[string]S, namedQueries map[string]map[string]FieldFilter) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if err := sc.mergeSourcesInternal(origin, sources); err != nil {
		return err
	}

	// Merge named queries (later origins override earlier ones at field level)
	for queryName, fieldFilters := range namedQueries {
		if sc.namedQueries[queryName] == nil {
			sc.namedQueries[queryName] = make(map[string]FieldFilter)
		}
		maps.Copy(sc.namedQueries[queryName], fieldFilters)
	}

	return nil
}

// mergeSourcesInternal performs the source merge. Must be called with lock held.
func (sc *GenericSourceCollection[S]) mergeSourcesInternal(origin string, sources map[string]S) error {
	// Find existing entry for this origin
	for i := range sc.entries {
		if sc.entries[i].origin == origin {
			sc.entries[i].sources = sources
			return nil
		}
	}

	// Origin not found, append it (dynamic registration)
	sc.entries = append(sc.entries, originEntry[S]{origin: origin, sources: sources})
	return nil
}

// GetNamedQuery returns a copy of a single named query by name.
// Returns (filters, true) if found, or (nil, false) if unknown.
// Slice values within each FieldFilter are cloned to prevent callers from
// accidentally mutating internal state.
func (sc *GenericSourceCollection[S]) GetNamedQuery(name string) (map[string]FieldFilter, bool) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	fieldFilters, ok := sc.namedQueries[name]
	if !ok {
		return nil, false
	}
	result := make(map[string]FieldFilter, len(fieldFilters))
	for field, ff := range fieldFilters {
		result[field] = DeepCopyFieldFilter(ff)
	}
	return result, true
}

// GetNamedQueries returns a deep copy of all merged named queries.
func (sc *GenericSourceCollection[S]) GetNamedQueries() map[string]map[string]FieldFilter {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	result := make(map[string]map[string]FieldFilter, len(sc.namedQueries))
	for queryName, fieldFilters := range sc.namedQueries {
		result[queryName] = make(map[string]FieldFilter, len(fieldFilters))
		for field, ff := range fieldFilters {
			result[queryName][field] = DeepCopyFieldFilter(ff)
		}
	}
	return result
}

// merged computes the merged view of all sources with field-level merging.
// Must be called with lock held.
func (sc *GenericSourceCollection[S]) merged() map[string]S {
	result := map[string]S{}

	for _, entry := range sc.entries {
		for id, source := range entry.sources {
			if existing, ok := result[id]; ok {
				// Field-level merge: existing is base, source is override
				result[id] = sc.mergeFunc(existing, source)
			} else {
				result[id] = source
			}
		}
	}

	// Apply defaults to all merged sources
	if sc.defaultsFunc != nil {
		for id, source := range result {
			result[id] = sc.defaultsFunc(source)
		}
	}

	return result
}

// AllSources returns all merged sources including Type and Properties.
// All sources are returned regardless of enabled status.
func (sc *GenericSourceCollection[S]) AllSources() map[string]S {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	result := make(map[string]S)
	maps.Copy(result, sc.merged())
	return result
}

// DeepCopyFieldFilter returns a copy of ff where slice values are cloned.
func DeepCopyFieldFilter(ff FieldFilter) FieldFilter {
	if vals, ok := ff.Value.([]any); ok {
		cp := make([]any, len(vals))
		copy(cp, vals)
		ff.Value = cp
	}
	return ff
}
