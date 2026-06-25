import * as React from 'react';
import { Navigate, Route, Routes } from 'react-router-dom';
import { NotFound } from 'mod-arch-shared';
import { useModularArchContext, DeploymentMode } from 'mod-arch-core';
import { NavDataItem } from '~/app/standalone/types';
import ModelRegistrySettingsRoutes from './pages/settings/ModelRegistrySettingsRoutes';
import ModelRegistryRoutes from './pages/modelRegistry/ModelRegistryRoutes';
import ModelCatalogRoutes from './pages/modelCatalog/ModelCatalogRoutes';
import McpCatalogRoutes from './pages/mcpCatalog/McpCatalogRoutes';
import AgentCatalogRoutes from './pages/agentCatalog/AgentCatalogRoutes';
import ModelCatalogSettingsRoutes from './pages/modelCatalogSettings/ModelCatalogSettingsRoutes';
import AgentCatalogSettingsRoutes from './pages/agentCatalogSettings/AgentCatalogSettingsRoutes';
import { modelCatalogUrl } from './routes/modelCatalog/catalogModel';
import { mcpCatalogUrl } from './routes/mcpCatalog/mcpCatalog';
import { agentCatalogUrl } from './routes/agentCatalog/agentCatalog';
import {
  catalogSettingsUrl,
  CATALOG_SETTINGS_PAGE_TITLE,
} from './routes/modelCatalogSettings/modelCatalogSettings';
import {
  agentCatalogSettingsUrl,
  AGENT_CATALOG_SETTINGS_PAGE_TITLE,
} from './routes/agentCatalogSettings/agentCatalogSettings';
import { modelRegistryUrl } from './pages/modelRegistry/screens/routeUtils';
import useUser from './hooks/useUser';

export const useAdminSettings = (): NavDataItem[] => {
  const { clusterAdmin } = useUser();
  const { config } = useModularArchContext();
  const { deploymentMode } = config;
  const isStandalone = deploymentMode === DeploymentMode.Standalone;
  const isFederated = deploymentMode === DeploymentMode.Federated;

  if (!clusterAdmin) {
    return [];
  }

  const settingsChildren = [{ label: 'Model registry settings', path: '/model-registry-settings' }];
  // Only show Model Catalog Settings in Standalone or Federated mode
  if (isStandalone || isFederated) {
    settingsChildren.push({ label: CATALOG_SETTINGS_PAGE_TITLE, path: catalogSettingsUrl() });
    settingsChildren.push({ label: AGENT_CATALOG_SETTINGS_PAGE_TITLE, path: agentCatalogSettingsUrl() });
  }

  return [
    {
      label: 'Settings',
      children: settingsChildren,
    },
  ];
};

export const useNavData = (): NavDataItem[] => {
  const { config } = useModularArchContext();
  const { deploymentMode } = config;
  const isStandalone = deploymentMode === DeploymentMode.Standalone;
  const isFederated = deploymentMode === DeploymentMode.Federated;

  const baseNavItems = [
    {
      label: 'Model Registry',
      path: modelRegistryUrl(),
    },
  ];

  if (isStandalone || isFederated) {
    baseNavItems.push(
      { label: 'Model Catalog', path: modelCatalogUrl() },
      { label: 'MCP Catalog', path: mcpCatalogUrl() },
      { label: 'Agent Catalog', path: agentCatalogUrl() },
    );
  }

  return [...baseNavItems, ...useAdminSettings()];
};

const AppRoutes: React.FC = () => {
  const { clusterAdmin } = useUser();
  const { config } = useModularArchContext();
  const { deploymentMode } = config;
  const isStandalone = deploymentMode === DeploymentMode.Standalone;
  const isFederated = deploymentMode === DeploymentMode.Federated;

  return (
    <Routes>
      <Route path="/" element={<Navigate to={modelRegistryUrl()} replace />} />
      <Route path={`${modelRegistryUrl()}/*`} element={<ModelRegistryRoutes />} />
      {(isStandalone || isFederated) && (
        <>
          <Route path={`${modelCatalogUrl()}/*`} element={<ModelCatalogRoutes />} />
          <Route path={`${mcpCatalogUrl()}/*`} element={<McpCatalogRoutes />} />
          <Route path={`${agentCatalogUrl()}/*`} element={<AgentCatalogRoutes />} />
          <Route path={`${catalogSettingsUrl()}/*`} element={<ModelCatalogSettingsRoutes />} />
          <Route path={`${agentCatalogSettingsUrl()}/*`} element={<AgentCatalogSettingsRoutes />} />
        </>
      )}
      <Route path="*" element={<NotFound />} />
      {/* TODO: [Conditional render] Follow up add testing and conditional rendering when in standalone mode */}
      {clusterAdmin && (
        <Route path="/model-registry-settings/*" element={<ModelRegistrySettingsRoutes />} />
      )}
    </Routes>
  );
};

export default AppRoutes;
