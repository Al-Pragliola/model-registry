package agentcatalog

import (
	"strings"
	"sync"

	"github.com/kubeflow/hub/catalog/internal/catalog/basecatalog"
)

type agentOriginEntry struct {
	origin  string
	sources map[string]basecatalog.AgentSource
}

// AgentSourceCollection manages agent catalog sources from multiple origins with priority-based merging.
type AgentSourceCollection struct {
	mu      sync.RWMutex
	entries []agentOriginEntry
}

func NewAgentSourceCollection(originOrder ...string) *AgentSourceCollection {
	entries := make([]agentOriginEntry, len(originOrder))
	for i, origin := range originOrder {
		entries[i] = agentOriginEntry{origin: origin, sources: nil}
	}
	return &AgentSourceCollection{
		entries: entries,
	}
}

func (sc *AgentSourceCollection) Merge(origin string, sources map[string]basecatalog.AgentSource) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	for i := range sc.entries {
		if sc.entries[i].origin == origin {
			sc.entries[i].sources = sources
			return nil
		}
	}

	sc.entries = append(sc.entries, agentOriginEntry{origin: origin, sources: sources})
	return nil
}

func (sc *AgentSourceCollection) merged() map[string]basecatalog.AgentSource {
	result := map[string]basecatalog.AgentSource{}

	for _, entry := range sc.entries {
		for id, source := range entry.sources {
			if existing, ok := result[id]; ok {
				result[id] = mergeAgentSources(existing, source)
			} else {
				result[id] = source
			}
		}
	}

	for id, source := range result {
		result[id] = applyAgentDefaults(source)
	}

	return result
}

func mergeAgentSources(base, override basecatalog.AgentSource) basecatalog.AgentSource {
	result := base

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
	if common.AssetType != nil {
		result.AssetType = common.AssetType
	}

	if override.IncludedAgents != nil {
		result.IncludedAgents = override.IncludedAgents
	}
	if override.ExcludedAgents != nil {
		result.ExcludedAgents = override.ExcludedAgents
	}

	return result
}

func applyAgentDefaults(source basecatalog.AgentSource) basecatalog.AgentSource {
	if source.Enabled == nil {
		source.Enabled = new(true)
	}
	if source.Labels == nil {
		source.Labels = []string{}
	}
	return source
}

func (sc *AgentSourceCollection) AllSources() map[string]basecatalog.AgentSource {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	return sc.merged()
}

// ByLabel returns enabled sources that have any of the labels provided.
// If a label is "null", every source without a label is returned.
func (sc *AgentSourceCollection) ByLabel(labels []string) []basecatalog.AgentSource {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	labelMap := make(map[string]struct{}, len(labels))
	for _, label := range labels {
		labelMap[strings.ToLower(label)] = struct{}{}
	}

	matches := map[string]basecatalog.AgentSource{}
	sources := sc.merged()

	if _, hasNull := labelMap["null"]; hasNull {
		for id, source := range sources {
			if source.Enabled == nil || !*source.Enabled {
				continue
			}
			if len(source.Labels) == 0 {
				matches[id] = source
			}
		}
	}

	for id, source := range sources {
		if source.Enabled == nil || !*source.Enabled {
			continue
		}
		for _, label := range source.Labels {
			if _, match := labelMap[strings.ToLower(label)]; match {
				matches[id] = source
				break
			}
		}
	}

	result := make([]basecatalog.AgentSource, 0, len(matches))
	for _, source := range matches {
		result = append(result, source)
	}
	return result
}
