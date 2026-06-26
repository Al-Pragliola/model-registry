import * as React from 'react';
import { Navigate, Route, Routes } from 'react-router-dom';
import { AgentCatalogContextProvider } from '~/app/context/agentCatalog/AgentCatalogContext';
import AgentCatalogCoreLoader from './AgentCatalogCoreLoader';
import AgentCatalog from './screens/AgentCatalog';
import AgentDetailsPage from './screens/AgentDetailsPage';

const AgentCatalogRoutes: React.FC = () => (
  <AgentCatalogContextProvider>
    <Routes>
      <Route path="/*" element={<AgentCatalogCoreLoader />}>
        <Route index element={<AgentCatalog />} />
        <Route path=":agentId" element={<AgentDetailsPage />} />
        <Route path="*" element={<Navigate to="." />} />
      </Route>
    </Routes>
  </AgentCatalogContextProvider>
);

export default AgentCatalogRoutes;
