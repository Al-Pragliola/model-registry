import { APIState, useAPIState } from 'mod-arch-core';
import React from 'react';
import {
  createCatalogSourceConfig,
  deleteCatalogSourceConfig,
  getCatalogSourceConfig,
  getCatalogSourceConfigs,
  updateCatalogSourceConfig,
  previewCatalogSource,
} from '~/app/api/agentCatalogSettings/service';
import { ModelCatalogSettingsAPIs } from '~/app/modelCatalogTypes';

export type AgentCatalogSettingsAPIState = APIState<ModelCatalogSettingsAPIs>;

const useAgentCatalogSettingsAPIState = (
  hostPath: string | null,
  queryParameters?: Record<string, unknown>,
): [apiState: AgentCatalogSettingsAPIState, refreshAPIState: () => void] => {
  const createAPI = React.useCallback(
    (path: string) => ({
      getCatalogSourceConfigs: getCatalogSourceConfigs(path, queryParameters),
      createCatalogSourceConfig: createCatalogSourceConfig(path, queryParameters),
      getCatalogSourceConfig: getCatalogSourceConfig(path, queryParameters),
      updateCatalogSourceConfig: updateCatalogSourceConfig(path, queryParameters),
      deleteCatalogSourceConfig: deleteCatalogSourceConfig(path, queryParameters),
      previewCatalogSource: previewCatalogSource(path, queryParameters),
    }),
    [queryParameters],
  );

  return useAPIState(hostPath, createAPI);
};

export default useAgentCatalogSettingsAPIState;
