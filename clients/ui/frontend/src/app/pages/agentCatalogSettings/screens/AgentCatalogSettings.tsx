import * as React from 'react';
import { Button, EmptyState, EmptyStateBody, EmptyStateVariant } from '@patternfly/react-core';
import { PlusCircleIcon } from '@patternfly/react-icons';
import { useNavigate } from 'react-router-dom';
import { ProjectObjectType, TitleWithIcon, ApplicationsPage } from 'mod-arch-shared';
import {
  AGENT_CATALOG_SETTINGS_PAGE_TITLE,
  AGENT_CATALOG_SETTINGS_DESCRIPTION,
  addAgentSourceUrl,
  ADD_AGENT_SOURCE_TITLE,
} from '~/app/routes/agentCatalogSettings/agentCatalogSettings';
import { AgentCatalogSettingsContext } from '~/app/context/agentCatalogSettings/AgentCatalogSettingsContext';
import AgentCatalogSourceConfigsTable from './AgentCatalogSourceConfigsTable';

const AgentCatalogSettings: React.FC = () => {
  const navigate = useNavigate();
  const {
    catalogSourceConfigs,
    catalogSourceConfigsLoaded,
    catalogSourceConfigsLoadError,
    apiState,
    refreshCatalogSourceConfigs,
  } = React.useContext(AgentCatalogSettingsContext);

  const configs = catalogSourceConfigs?.catalogs || [];
  const isEmpty = catalogSourceConfigsLoaded && configs.length === 0;

  const handleDeleteSource = React.useCallback(
    async (sourceId: string): Promise<void> => {
      if (!apiState.apiAvailable) {
        throw new Error('API not available');
      }
      await apiState.api.deleteCatalogSourceConfig({}, sourceId);
      refreshCatalogSourceConfigs();
    },
    [apiState.api, apiState.apiAvailable, refreshCatalogSourceConfigs],
  );

  return (
    <ApplicationsPage
      title={
        <TitleWithIcon
          title={AGENT_CATALOG_SETTINGS_PAGE_TITLE}
          objectType={ProjectObjectType.modelCatalog}
        />
      }
      description={AGENT_CATALOG_SETTINGS_DESCRIPTION}
      empty={isEmpty}
      emptyStatePage={
        <EmptyState
          headingLevel="h5"
          icon={PlusCircleIcon}
          titleText="No agent catalog sources"
          variant={EmptyStateVariant.lg}
          data-testid="agent-catalog-settings-empty-state"
        >
          <EmptyStateBody>
            No agent catalog sources have been configured. Add a source to get started.
          </EmptyStateBody>
          <Button
            variant="primary"
            onClick={() => navigate(addAgentSourceUrl())}
            data-testid="add-agent-source-button-empty"
          >
            {ADD_AGENT_SOURCE_TITLE}
          </Button>
        </EmptyState>
      }
      loaded={catalogSourceConfigsLoaded}
      loadError={catalogSourceConfigsLoadError}
      errorMessage="Unable to load agent catalog source configurations."
      provideChildrenPadding
    >
      <AgentCatalogSourceConfigsTable
        catalogSourceConfigs={configs}
        onAddSource={() => navigate(addAgentSourceUrl())}
        onDeleteSource={handleDeleteSource}
      />
    </ApplicationsPage>
  );
};

export default AgentCatalogSettings;
