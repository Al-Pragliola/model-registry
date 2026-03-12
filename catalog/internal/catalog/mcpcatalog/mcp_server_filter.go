package mcpcatalog

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/golang/glog"
	"github.com/kubeflow/model-registry/catalog/internal/catalog/basecatalog"
)

// MCPServerFilter encapsulates include/exclude pattern matching for MCP server names.
type MCPServerFilter struct {
	included []*mcpCompiledPattern
	excluded []*mcpCompiledPattern
}

type mcpCompiledPattern struct {
	raw string
	re  *regexp.Regexp
}

func newMCPCompiledPattern(field string, idx int, raw string) (*mcpCompiledPattern, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return nil, fmt.Errorf("%s[%d]: pattern cannot be empty", field, idx)
	}

	// Convert a simple glob (only supporting '*') into a regexp.
	var b strings.Builder
	b.WriteString("(?i)^")
	for _, r := range value {
		if r == '*' {
			b.WriteString(".*")
			continue
		}
		b.WriteString(regexp.QuoteMeta(string(r)))
	}
	b.WriteString("$")

	re, err := regexp.Compile(b.String())
	if err != nil {
		return nil, fmt.Errorf("%s[%d]: invalid pattern %q: %w", field, idx, value, err)
	}

	return &mcpCompiledPattern{
		raw: value,
		re:  re,
	}, nil
}

func compileMCPPatterns(field string, patterns []string) ([]*mcpCompiledPattern, error) {
	if len(patterns) == 0 {
		return nil, nil
	}

	compiled := make([]*mcpCompiledPattern, 0, len(patterns))
	for i, pattern := range patterns {
		cp, err := newMCPCompiledPattern(field, i, pattern)
		if err != nil {
			return nil, err
		}
		compiled = append(compiled, cp)
	}
	return compiled, nil
}

// ValidateMCPServerFilters validates that the includedServers and excludedServers patterns
// are valid (non-empty, compilable). This is useful for early validation
// at configuration load time without constructing the full MCPServerFilter.
func ValidateMCPServerFilters(included, excluded []string) error {
	detectConflictingMCPPatterns(included, excluded)

	if _, err := compileMCPPatterns("includedServers", included); err != nil {
		return err
	}

	if _, err := compileMCPPatterns("excludedServers", excluded); err != nil {
		return err
	}

	return nil
}

// NewMCPServerFilter builds an MCPServerFilter from the provided include/exclude pattern lists.
// Returns nil if both lists are empty (no filtering needed).
func NewMCPServerFilter(included, excluded []string) (*MCPServerFilter, error) {
	if err := ValidateMCPServerFilters(included, excluded); err != nil {
		return nil, err
	}

	inc, err := compileMCPPatterns("includedServers", included)
	if err != nil {
		return nil, err
	}

	exc, err := compileMCPPatterns("excludedServers", excluded)
	if err != nil {
		return nil, err
	}

	if len(inc) == 0 && len(exc) == 0 {
		return nil, nil
	}

	return &MCPServerFilter{
		included: inc,
		excluded: exc,
	}, nil
}

func detectConflictingMCPPatterns(included, excluded []string) {
	if len(included) == 0 || len(excluded) == 0 {
		return
	}

	includedIdx := make(map[string]int, len(included))
	for i, pattern := range included {
		value := strings.TrimSpace(pattern)
		includedIdx[value] = i
	}

	for _, pattern := range excluded {
		value := strings.TrimSpace(pattern)
		if _, exists := includedIdx[value]; exists {
			glog.Errorf("pattern %q is defined in both includedServers and excludedServers, skipping", value)
		}
	}
}

// Allows returns true if the provided server name passes the include/exclude rules.
// Exclusions take precedence over inclusions.
// A nil filter allows all server names.
func (f *MCPServerFilter) Allows(name string) bool {
	if f == nil {
		return true
	}

	// Check exclusions first — they take precedence
	for _, pattern := range f.excluded {
		if pattern.re.MatchString(name) {
			return false
		}
	}

	// If no inclusions, allow everything that wasn't excluded
	if len(f.included) == 0 {
		return true
	}

	// Check if the name matches any inclusion pattern
	for _, pattern := range f.included {
		if pattern.re.MatchString(name) {
			return true
		}
	}

	return false
}

// NewMCPServerFilterFromSource builds an MCPServerFilter from the source-level configuration.
func NewMCPServerFilterFromSource(source *basecatalog.MCPSource) (*MCPServerFilter, error) {
	if source == nil {
		return nil, fmt.Errorf("source cannot be nil when building MCP server filters")
	}

	filter, err := NewMCPServerFilter(source.IncludedServers, source.ExcludedServers)
	if err != nil {
		return nil, fmt.Errorf("invalid include/exclude configuration for MCP source %s: %w", source.ID, err)
	}

	return filter, nil
}
