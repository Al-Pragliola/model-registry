import type { ModelCatalogAPIState } from '~/app/hooks/modelCatalog/useModelCatalogAPIState';
import type { CatalogContextValue } from '~/app/context/catalogContext/createCatalogContext';
import type {
  AgentCatalogFilterOptionsList,
  AgentCatalogFiltersState,
} from '~/app/pages/agentCatalog/types/agentCatalogFilterOptions';

export type AgentCatalogPaginationState = {
  page: number;
  pageSize: number;
  totalItems: number;
};

export type AgentCatalogExtension = {
  filters: AgentCatalogFiltersState;
  setFilters: (
    filters: AgentCatalogFiltersState | ((prev: AgentCatalogFiltersState) => AgentCatalogFiltersState),
  ) => void;
  searchQuery: string;
  setSearchQuery: (query: string) => void;
  pagination: AgentCatalogPaginationState;
  setPage: (page: number) => void;
  setPageSize: (pageSize: number) => void;
  setTotalItems: (totalItems: number) => void;
  clearAllFilters: () => void;
  agentApiState: ModelCatalogAPIState;
};

export type AgentCatalogContextType = CatalogContextValue<AgentCatalogFilterOptionsList> &
  AgentCatalogExtension;
