import { FetchState, FetchStateCallbackPromise, NotReadyError, useFetchState } from 'mod-arch-core';
import React from 'react';
import { CatalogSourceConfig } from '~/app/modelCatalogTypes';
import { AgentCatalogSettingsContext } from '~/app/context/agentCatalogSettings/AgentCatalogSettingsContext';

type State = CatalogSourceConfig | null;

export const useAgentCatalogSourceConfigBySourceId = (sourceId: string): FetchState<State> => {
  const { apiState } = React.useContext(AgentCatalogSettingsContext);
  const call = React.useCallback<FetchStateCallbackPromise<State>>(
    (opts) => {
      if (!apiState.apiAvailable) {
        return Promise.reject(new Error('API not yet available'));
      }
      if (!sourceId) {
        return Promise.reject(new NotReadyError('No source id'));
      }

      return apiState.api.getCatalogSourceConfig(opts, sourceId);
    },
    [apiState, sourceId],
  );
  return useFetchState(call, null, { initialPromisePurity: true });
};
