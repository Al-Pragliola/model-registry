package agentcatalog

import (
	"github.com/kubeflow/hub/catalog/internal/catalog/basecatalog"
)

type AgentFilter = basecatalog.NameFilter

func ValidateAgentSourceFilters(included, excluded []string) error {
	return basecatalog.ValidatePatterns("includedAgents", included, "excludedAgents", excluded)
}

func NewAgentFilter(included, excluded []string) (*AgentFilter, error) {
	return basecatalog.NewNameFilter("includedAgents", included, "excludedAgents", excluded)
}
