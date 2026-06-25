import * as React from 'react';
import { useQueryParamNamespaces } from 'mod-arch-core';
import useAgentCatalogSettingsAPIState, {
  AgentCatalogSettingsAPIState,
} from '~/app/hooks/agentCatalogSettings/useAgentCatalogSettingsAPIState';
import { useAgentCatalogSourceConfigs } from '~/app/hooks/agentCatalogSettings/useAgentCatalogSourceConfigs';
import type { CatalogSourceList } from '~/app/shared/types/catalogTypes';
import type { CatalogSourceConfigList } from '~/app/modelCatalogTypes';
import { BFF_API_VERSION, URL_PREFIX } from '~/app/utilities/const';
import useModelCatalogAPIState from '~/app/hooks/modelCatalog/useModelCatalogAPIState';
import { useCatalogSourcesWithPolling } from '~/app/hooks/modelCatalogSettings/useCatalogSourcesWithPolling';

export type AgentCatalogSettingsContextType = {
  apiState: AgentCatalogSettingsAPIState;
  refreshAPIState: () => void;
  catalogSourceConfigs: CatalogSourceConfigList | null;
  catalogSourceConfigsLoaded: boolean;
  catalogSourceConfigsLoadError?: Error;
  refreshCatalogSourceConfigs: () => void;
  catalogSources: CatalogSourceList | null;
  catalogSourcesLoaded: boolean;
  catalogSourcesLoadError?: Error;
  refreshCatalogSources: () => void;
};

type AgentCatalogSettingsContextProviderProps = {
  children: React.ReactNode;
};

export const AgentCatalogSettingsContext = React.createContext<AgentCatalogSettingsContextType>({
  // eslint-disable-next-line @typescript-eslint/consistent-type-assertions
  apiState: { apiAvailable: false, api: null as unknown as AgentCatalogSettingsAPIState['api'] },
  refreshAPIState: () => undefined,
  catalogSourceConfigs: null,
  catalogSourceConfigsLoaded: false,
  catalogSourceConfigsLoadError: undefined,
  refreshCatalogSourceConfigs: () => undefined,
  catalogSources: null,
  catalogSourcesLoaded: false,
  catalogSourcesLoadError: undefined,
  refreshCatalogSources: () => undefined,
});

export const AgentCatalogSettingsContextProvider: React.FC<
  AgentCatalogSettingsContextProviderProps
> = ({ children }) => {
  const hostPath = `${URL_PREFIX}/api/${BFF_API_VERSION}/settings/agent_catalog`;
  const catalogHostPath = `${URL_PREFIX}/api/${BFF_API_VERSION}/agent_catalog`;
  const queryParams = useQueryParamNamespaces();
  const [apiState, refreshAPIState] = useAgentCatalogSettingsAPIState(hostPath, queryParams);
  const [catalogAPIState] = useModelCatalogAPIState(catalogHostPath, queryParams);
  const [
    catalogSourceConfigs,
    catalogSourceConfigsLoaded,
    catalogSourceConfigsLoadError,
    refreshCatalogSourceConfigs,
  ] = useAgentCatalogSourceConfigs(apiState);

  // Fetch catalog sources with polling for status updates
  const [catalogSources, catalogSourcesLoaded, catalogSourcesLoadError, refreshCatalogSources] =
    useCatalogSourcesWithPolling(catalogAPIState);

  const contextValue = React.useMemo(
    () => ({
      apiState,
      refreshAPIState,
      catalogSourceConfigs,
      catalogSourceConfigsLoaded,
      catalogSourceConfigsLoadError,
      refreshCatalogSourceConfigs,
      catalogSources,
      catalogSourcesLoaded,
      catalogSourcesLoadError,
      refreshCatalogSources,
    }),
    [
      apiState,
      refreshAPIState,
      catalogSourceConfigs,
      catalogSourceConfigsLoaded,
      catalogSourceConfigsLoadError,
      refreshCatalogSourceConfigs,
      catalogSources,
      catalogSourcesLoaded,
      catalogSourcesLoadError,
      refreshCatalogSources,
    ],
  );

  return (
    <AgentCatalogSettingsContext.Provider value={contextValue}>
      {children}
    </AgentCatalogSettingsContext.Provider>
  );
};
