package mcpcatalog

import (
	"testing"

	"github.com/kubeflow/model-registry/catalog/internal/catalog/basecatalog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMCPServerFilterAllows(t *testing.T) {
	filter, err := NewMCPServerFilter([]string{"github-*", "slack-*"}, []string{"*-deprecated", "*-experimental"})
	require.NoError(t, err)

	// Included servers should pass
	assert.True(t, filter.Allows("github-mcp"))
	assert.True(t, filter.Allows("github-actions-mcp"))
	assert.True(t, filter.Allows("slack-mcp"))
	assert.True(t, filter.Allows("slack-notifications"))

	// Non-included servers should be rejected
	assert.False(t, filter.Allows("jira-mcp"))
	assert.False(t, filter.Allows("random-server"))

	// Excluded servers should be rejected even if they match inclusion
	assert.False(t, filter.Allows("github-deprecated"))
	assert.False(t, filter.Allows("slack-experimental"))
}

func TestMCPServerFilterCaseInsensitive(t *testing.T) {
	filter, err := NewMCPServerFilter([]string{"GitHub-*"}, []string{"*-Alpha"})
	require.NoError(t, err)

	assert.True(t, filter.Allows("github-mcp"))
	assert.True(t, filter.Allows("GITHUB-MCP"))
	assert.True(t, filter.Allows("GitHub-Actions"))
	assert.False(t, filter.Allows("github-alpha"))
	assert.False(t, filter.Allows("GITHUB-ALPHA"))
}

func TestMCPServerFilterExclusionsPrecedence(t *testing.T) {
	// Exclusions should take precedence over inclusions
	filter, err := NewMCPServerFilter([]string{"*"}, []string{"*-beta", "*-alpha"})
	require.NoError(t, err)

	assert.True(t, filter.Allows("github-mcp"))
	assert.True(t, filter.Allows("slack-mcp"))
	assert.False(t, filter.Allows("github-beta"))
	assert.False(t, filter.Allows("slack-alpha"))
}

func TestMCPServerFilterEmptyAllowsAll(t *testing.T) {
	// No filters should allow everything
	filter, err := NewMCPServerFilter(nil, nil)
	require.NoError(t, err)
	assert.Nil(t, filter)

	// A nil filter should allow all
	var nilFilter *MCPServerFilter
	assert.True(t, nilFilter.Allows("anything"))
	assert.True(t, nilFilter.Allows("github-mcp"))
}

func TestMCPServerFilterOnlyExclusions(t *testing.T) {
	// Only exclusions — everything allowed except excluded patterns
	filter, err := NewMCPServerFilter(nil, []string{"*-deprecated", "*-experimental", "*-alpha"})
	require.NoError(t, err)

	assert.True(t, filter.Allows("github-mcp"))
	assert.True(t, filter.Allows("slack-mcp"))
	assert.True(t, filter.Allows("jira-server"))
	assert.False(t, filter.Allows("github-deprecated"))
	assert.False(t, filter.Allows("slack-experimental"))
	assert.False(t, filter.Allows("jira-alpha"))
}

func TestMCPServerFilterOnlyInclusions(t *testing.T) {
	// Only inclusions — only matching servers allowed
	filter, err := NewMCPServerFilter([]string{"github-*", "slack-*"}, nil)
	require.NoError(t, err)

	assert.True(t, filter.Allows("github-mcp"))
	assert.True(t, filter.Allows("slack-mcp"))
	assert.False(t, filter.Allows("jira-mcp"))
	assert.False(t, filter.Allows("random-server"))
}

func TestMCPServerFilterPrefixPatterns(t *testing.T) {
	filter, err := NewMCPServerFilter([]string{"github-*"}, nil)
	require.NoError(t, err)

	assert.True(t, filter.Allows("github-mcp"))
	assert.True(t, filter.Allows("github-actions-mcp"))
	assert.True(t, filter.Allows("github-"))
	assert.False(t, filter.Allows("gitlab-mcp"))
}

func TestMCPServerFilterSuffixPatterns(t *testing.T) {
	filter, err := NewMCPServerFilter([]string{"*-mcp"}, nil)
	require.NoError(t, err)

	assert.True(t, filter.Allows("github-mcp"))
	assert.True(t, filter.Allows("slack-mcp"))
	assert.True(t, filter.Allows("jira-mcp"))
	assert.False(t, filter.Allows("github-server"))
	assert.False(t, filter.Allows("mcp-github"))
}

func TestMCPServerFilterMiddlePatterns(t *testing.T) {
	filter, err := NewMCPServerFilter([]string{"*hub*"}, nil)
	require.NoError(t, err)

	assert.True(t, filter.Allows("github-mcp"))
	assert.True(t, filter.Allows("hub-server"))
	assert.True(t, filter.Allows("my-hub-thing"))
	assert.False(t, filter.Allows("slack-mcp"))
}

func TestMCPServerFilterExactMatch(t *testing.T) {
	filter, err := NewMCPServerFilter([]string{"github-mcp"}, nil)
	require.NoError(t, err)

	assert.True(t, filter.Allows("github-mcp"))
	assert.True(t, filter.Allows("GitHub-MCP")) // case-insensitive
	assert.False(t, filter.Allows("github-mcp-v2"))
	assert.False(t, filter.Allows("my-github-mcp"))
}

func TestMCPServerFilterMultiplePatterns(t *testing.T) {
	filter, err := NewMCPServerFilter(
		[]string{"github-*", "slack-*", "jira-*"},
		[]string{"*-deprecated", "*-experimental", "*-alpha"},
	)
	require.NoError(t, err)

	// Matches inclusion patterns
	assert.True(t, filter.Allows("github-mcp"))
	assert.True(t, filter.Allows("slack-notifications"))
	assert.True(t, filter.Allows("jira-issues"))

	// Doesn't match any inclusion pattern
	assert.False(t, filter.Allows("gitlab-mcp"))
	assert.False(t, filter.Allows("confluence-mcp"))

	// Matches exclusion patterns (takes precedence)
	assert.False(t, filter.Allows("github-deprecated"))
	assert.False(t, filter.Allows("slack-experimental"))
	assert.False(t, filter.Allows("jira-alpha"))
}

func TestMCPServerFilterValidation(t *testing.T) {
	t.Run("empty pattern rejected", func(t *testing.T) {
		_, err := NewMCPServerFilter([]string{""}, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "pattern cannot be empty")
	})

	t.Run("whitespace-only pattern rejected", func(t *testing.T) {
		_, err := NewMCPServerFilter([]string{"   "}, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "pattern cannot be empty")
	})

	t.Run("valid patterns accepted", func(t *testing.T) {
		err := ValidateMCPServerFilters([]string{"github-*", "*-mcp"}, []string{"*-beta"})
		assert.NoError(t, err)
	})

	t.Run("no filters valid", func(t *testing.T) {
		err := ValidateMCPServerFilters(nil, nil)
		assert.NoError(t, err)
	})

	t.Run("conflicting patterns logged but not rejected", func(t *testing.T) {
		_, err := NewMCPServerFilter([]string{"github-*"}, []string{"github-*"})
		require.NoError(t, err) // conflicts are logged, not rejected
	})
}

func TestMCPServerFilterWildcardAll(t *testing.T) {
	filter, err := NewMCPServerFilter([]string{"*"}, nil)
	require.NoError(t, err)

	assert.True(t, filter.Allows("anything"))
	assert.True(t, filter.Allows("github-mcp"))
	assert.True(t, filter.Allows(""))
}

func TestNewMCPServerFilterFromSource(t *testing.T) {
	t.Run("source with both filters", func(t *testing.T) {
		source := &basecatalog.MCPSource{
			ID:              "test-source",
			IncludedServers: []string{"github-*", "slack-*"},
			ExcludedServers: []string{"*-deprecated"},
		}

		filter, err := NewMCPServerFilterFromSource(source)
		require.NoError(t, err)
		require.NotNil(t, filter)

		assert.True(t, filter.Allows("github-mcp"))
		assert.True(t, filter.Allows("slack-mcp"))
		assert.False(t, filter.Allows("jira-mcp"))
		assert.False(t, filter.Allows("github-deprecated"))
	})

	t.Run("source with no filters", func(t *testing.T) {
		source := &basecatalog.MCPSource{
			ID: "test-source",
		}

		filter, err := NewMCPServerFilterFromSource(source)
		require.NoError(t, err)
		assert.Nil(t, filter)
	})

	t.Run("nil source returns error", func(t *testing.T) {
		_, err := NewMCPServerFilterFromSource(nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "source cannot be nil")
	})

	t.Run("source with invalid pattern", func(t *testing.T) {
		source := &basecatalog.MCPSource{
			ID:              "test-source",
			IncludedServers: []string{""},
		}

		_, err := NewMCPServerFilterFromSource(source)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "pattern cannot be empty")
	})
}
