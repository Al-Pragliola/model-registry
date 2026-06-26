import type { AgentCatalogFiltersState } from '~/app/pages/agentCatalog/types/agentCatalogFilterOptions';
import {
  BACKEND_TO_FRONTEND_AGENT_FILTER_KEY,
  AGENT_FILTER_KEYS,
} from '~/app/pages/agentCatalog/const';
import { stringFiltersToFilterQuery } from '~/app/shared/components/catalog';

export const hasAgentFiltersApplied = (
  filters: AgentCatalogFiltersState,
  searchQuery: string,
): boolean => {
  if (searchQuery && searchQuery.trim().length > 0) {
    return true;
  }
  for (const key of AGENT_FILTER_KEYS) {
    const value = filters[key];
    if (Array.isArray(value) && value.length > 0) {
      return true;
    }
  }
  return false;
};

const FRONTEND_TO_BACKEND_FILTER_KEY: Record<string, string> = Object.fromEntries(
  Object.entries(BACKEND_TO_FRONTEND_AGENT_FILTER_KEY).map(([backend, frontend]) => [
    frontend,
    backend,
  ]),
);

export function agentFiltersToFilterQuery(filters: AgentCatalogFiltersState): string {
  return stringFiltersToFilterQuery(filters, FRONTEND_TO_BACKEND_FILTER_KEY);
}
