package basecatalog

import (
	"fmt"
	"path/filepath"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/golang/glog"
)

// DelegateLoaderConfig defines the pluggable behaviors for a DelegateLoader.
// Each catalog type (model, MCP) provides its own implementation.
type DelegateLoaderConfig[S Source] struct {
	// Name is a human-readable label for log messages (e.g. "model", "MCP").
	Name string

	// ExtractSources returns the sources of this type from a parsed SourceConfig.
	ExtractSources func(config *SourceConfig, origin string) (map[string]S, error)

	// ExtractNamedQueries returns named queries from a parsed config, or nil.
	ExtractNamedQueries func(config *SourceConfig) map[string]map[string]FieldFilter

	// SourceCollection is the collection to merge sources into.
	SourceCollection *GenericSourceCollection[S]
}

// DelegateLoader provides the common lifecycle operations shared by all
// catalog-type loaders (model catalogs, MCP catalogs, etc.).
//
// It handles:
//   - Parsing config files and merging sources into a GenericSourceCollection
//   - Coordinating with LoaderState for leader/write tracking
//
// Catalog-specific operations (loading from providers, saving to DB, orphan
// cleanup) remain in the domain-specific loader since they differ too much
// to generalize without losing clarity.
type DelegateLoader[S Source] struct {
	state  LoaderState
	config DelegateLoaderConfig[S]
}

// NewDelegateLoader creates a new DelegateLoader with the given state and config.
func NewDelegateLoader[S Source](state LoaderState, config DelegateLoaderConfig[S]) *DelegateLoader[S] {
	return &DelegateLoader[S]{
		state:  state,
		config: config,
	}
}

// ParseAllConfigs parses all config files and merges sources into the collection.
// This is the common initialization step for all loader types.
func (dl *DelegateLoader[S]) ParseAllConfigs() error {
	glog.Infof("Initializing %s loader - parsing configs", dl.config.Name)

	for _, path := range dl.state.Paths() {
		if err := dl.parseAndMerge(path); err != nil {
			return fmt.Errorf("failed to parse %s config %s: %w", dl.config.Name, path, err)
		}
	}

	glog.Infof("%s loader config parsing complete", dl.config.Name)
	return nil
}

// ReloadParsing re-parses all config files into in-memory collections.
// Called by the unified loader before computing combined source IDs.
func (dl *DelegateLoader[S]) ReloadParsing() {
	for _, path := range dl.state.Paths() {
		if err := dl.parseAndMerge(path); err != nil {
			glog.Errorf("unable to reload %s sources from %s: %v", dl.config.Name, path, err)
		}
	}
}

// CollectSourceIDs returns the set of all source IDs in the collection.
func (dl *DelegateLoader[S]) CollectSourceIDs() mapset.Set[string] {
	ids := mapset.NewSet[string]()
	for id := range dl.config.SourceCollection.AllSources() {
		ids.Add(id)
	}
	return ids
}

// parseAndMerge parses a config file and merges its sources into the collection.
func (dl *DelegateLoader[S]) parseAndMerge(path string) error {
	path, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for %s: %v", path, err)
	}

	config, err := ReadSourceConfig(path)
	if err != nil {
		return err
	}

	// Extract sources using the pluggable function
	sources, err := dl.config.ExtractSources(config, path)
	if err != nil {
		return err
	}

	// Merge using named queries if available
	if dl.config.ExtractNamedQueries != nil {
		namedQueries := dl.config.ExtractNamedQueries(config)
		if namedQueries != nil {
			return dl.config.SourceCollection.MergeWithNamedQueries(path, sources, namedQueries)
		}
	}

	return dl.config.SourceCollection.Merge(path, sources)
}
