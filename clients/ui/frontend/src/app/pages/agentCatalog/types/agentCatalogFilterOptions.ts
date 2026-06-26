import type { CatalogFilterStringOption } from '~/app/shared/components/catalog';

export type AgentFilterCategoryKey = 'framework' | 'agentType' | 'labels';

export type AgentCatalogFiltersState = {
  [K in AgentFilterCategoryKey]?: string[];
};

export type AgentCatalogFilterOptions = {
  [key in AgentFilterCategoryKey]?: CatalogFilterStringOption;
};

export type AgentCatalogFilterOptionsList = {
  filters?: AgentCatalogFilterOptions;
};
