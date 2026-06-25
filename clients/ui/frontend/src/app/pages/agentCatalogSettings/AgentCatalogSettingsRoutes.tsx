import * as React from 'react';
import { Navigate, Route, Routes } from 'react-router-dom';
import { AgentCatalogSettingsContextProvider } from '~/app/context/agentCatalogSettings/AgentCatalogSettingsContext';
import AgentCatalogSettings from '~/app/pages/agentCatalogSettings/screens/AgentCatalogSettings';
import ManageAgentSourcePage from '~/app/pages/agentCatalogSettings/screens/ManageAgentSourcePage';

const AgentCatalogSettingsRoutes: React.FC = () => (
  <AgentCatalogSettingsContextProvider>
    <Routes>
      <Route path="/" element={<AgentCatalogSettings />} />
      <Route path="add-source" element={<ManageAgentSourcePage />} />
      <Route path="manage-source/:catalogSourceId" element={<ManageAgentSourcePage />} />
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  </AgentCatalogSettingsContextProvider>
);

export default AgentCatalogSettingsRoutes;
