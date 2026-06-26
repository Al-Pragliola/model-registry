import * as React from 'react';
import { useQueryParamNamespaces } from 'mod-arch-core';
import { BFF_API_VERSION, URL_PREFIX } from '~/app/utilities/const';
import {
  createCatalogContext,
  CatalogCommonData,
  CatalogProviderState,
} from '~/app/context/catalogContext/createCatalogContext';
import useModelCatalogAPIState from '~/app/hooks/modelCatalog/useModelCatalogAPIState';
import { useCatalogSources } from '~/app/hooks/modelCatalog/useCatalogSources';
import { useCatalogLabels } from '~/app/hooks/modelCatalog/useCatalogLabels';
import { useAgentFilterOptionListWithAPI } from '~/app/hooks/agentCatalog/useAgentFilterOptionList';
import type {
  AgentCatalogFiltersState,
  AgentCatalogFilterOptionsList,
} from '~/app/pages/agentCatalog/types/agentCatalogFilterOptions';
import { useAgentUrlSync } from '~/app/pages/agentCatalog/hooks/useAgentUrlSync';
import type { AgentCatalogExtension, AgentCatalogPaginationState } from './types';

export type {
  AgentCatalogContextType,
  AgentCatalogExtension,
  AgentCatalogPaginationState,
} from './types';
export type { AgentCatalogFiltersState } from '~/app/pages/agentCatalog/types/agentCatalogFilterOptions';

const MODEL_CATALOG_PATH = `${URL_PREFIX}/api/${BFF_API_VERSION}/model_catalog`;
const AGENT_CATALOG_PATH = `${URL_PREFIX}/api/${BFF_API_VERSION}/agent_catalog`;

const defaultPagination: AgentCatalogPaginationState = {
  page: 1,
  pageSize: 10,
  totalItems: 0,
};

function useAgentCatalogSetup(providerState: CatalogProviderState) {
  const queryParams = useQueryParamNamespaces();
  const [apiStateModelCatalog] = useModelCatalogAPIState(MODEL_CATALOG_PATH, queryParams);
  const [apiStateAgentCatalog] = useModelCatalogAPIState(AGENT_CATALOG_PATH, queryParams);

  const agentListParams = React.useMemo(() => ({ assetType: 'agents' as const }), []);
  const [catalogSources, catalogSourcesLoaded, catalogSourcesLoadError] = useCatalogSources(
    apiStateModelCatalog,
    agentListParams,
  );
  const [catalogLabels, catalogLabelsLoaded, catalogLabelsLoadError] = useCatalogLabels(
    apiStateModelCatalog,
    agentListParams,
  );
  const [filterOptions, filterOptionsLoaded, filterOptionsLoadError] =
    useAgentFilterOptionListWithAPI(apiStateAgentCatalog);

  const { initialState, syncToUrl } = useAgentUrlSync();

  const [filters, setFilters] = React.useState<AgentCatalogFiltersState>(initialState.filters);
  const [searchQuery, setSearchQuery] = React.useState(initialState.searchQuery);
  const [pagination, setPaginationState] =
    React.useState<AgentCatalogPaginationState>(defaultPagination);

  const { setSelectedSourceLabel } = providerState;

  React.useEffect(() => {
    setSelectedSourceLabel(initialState.selectedSourceLabel);
  }, [setSelectedSourceLabel, initialState.selectedSourceLabel]);

  React.useEffect(() => {
    syncToUrl({
      searchQuery,
      filters,
      selectedSourceLabel: providerState.selectedSourceLabel,
    });
  }, [searchQuery, filters, providerState.selectedSourceLabel, syncToUrl]);

  const setPage = React.useCallback((page: number) => {
    setPaginationState((prev) => ({ ...prev, page }));
  }, []);

  const setPageSize = React.useCallback((pageSize: number) => {
    setPaginationState((prev) => ({ ...prev, pageSize, page: 1 }));
  }, []);

  const setTotalItems = React.useCallback((totalItems: number) => {
    setPaginationState((prev) => ({ ...prev, totalItems }));
  }, []);

  const clearAllFilters = React.useCallback(() => {
    setSearchQuery('');
    setFilters({});
  }, []);

  const catalogData = React.useMemo<CatalogCommonData<AgentCatalogFilterOptionsList>>(
    () => ({
      catalogSources,
      catalogSourcesLoaded,
      catalogSourcesLoadError,
      catalogLabels,
      catalogLabelsLoaded,
      catalogLabelsLoadError,
      filterOptions,
      filterOptionsLoaded,
      filterOptionsLoadError,
    }),
    [
      catalogSources,
      catalogSourcesLoaded,
      catalogSourcesLoadError,
      catalogLabels,
      catalogLabelsLoaded,
      catalogLabelsLoadError,
      filterOptions,
      filterOptionsLoaded,
      filterOptionsLoadError,
    ],
  );

  const extension = React.useMemo(
    () => ({
      filters,
      setFilters,
      searchQuery,
      setSearchQuery,
      pagination,
      setPage,
      setPageSize,
      setTotalItems,
      clearAllFilters,
      agentApiState: apiStateAgentCatalog,
    }),
    [
      apiStateAgentCatalog,
      filters,
      searchQuery,
      pagination,
      setPage,
      setPageSize,
      setTotalItems,
      clearAllFilters,
    ],
  );

  return { catalogData, extension };
}

const {
  Context: AgentCatalogContext,
  Provider: AgentCatalogContextProvider,
  useContext: useAgentCatalogContext,
} = createCatalogContext<AgentCatalogFilterOptionsList, AgentCatalogExtension>({
  displayName: 'AgentCatalogContextProvider',
  useSetup: useAgentCatalogSetup,
});

export { AgentCatalogContext, AgentCatalogContextProvider, useAgentCatalogContext };
