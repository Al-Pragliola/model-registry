import * as React from 'react';
import { useParams, Link } from 'react-router-dom';
import {
  Breadcrumb,
  BreadcrumbItem,
  Drawer,
  DrawerContent,
  DrawerContentBody,
  DrawerPanelContent,
} from '@patternfly/react-core';
import { ApplicationsPage } from 'mod-arch-shared';
import {
  AGENT_CATALOG_SETTINGS_PAGE_TITLE,
  ADD_AGENT_SOURCE_TITLE,
  ADD_AGENT_SOURCE_DESCRIPTION,
  MANAGE_AGENT_SOURCE_TITLE,
  MANAGE_AGENT_SOURCE_DESCRIPTION,
  agentCatalogSettingsUrl,
} from '~/app/routes/agentCatalogSettings/agentCatalogSettings';
import ManageAgentSourceForm from '~/app/pages/agentCatalogSettings/components/ManageAgentSourceForm';
import { ExpectedAgentYamlFormatDrawerPanel } from '~/app/pages/agentCatalogSettings/components/ExpectedAgentYamlFormatDrawer';
import { useAgentCatalogSourceConfigBySourceId } from '~/app/hooks/agentCatalogSettings/useAgentCatalogSourceConfigBySourceId';

const ManageAgentSourcePage: React.FC = () => {
  const { catalogSourceId } = useParams<{ catalogSourceId?: string }>();
  const isAddMode = !catalogSourceId;
  const pageTitle = isAddMode ? ADD_AGENT_SOURCE_TITLE : MANAGE_AGENT_SOURCE_TITLE;
  const breadcrumbLabel = isAddMode ? ADD_AGENT_SOURCE_TITLE : MANAGE_AGENT_SOURCE_TITLE;
  const description = isAddMode ? ADD_AGENT_SOURCE_DESCRIPTION : MANAGE_AGENT_SOURCE_DESCRIPTION;

  const state = useAgentCatalogSourceConfigBySourceId(catalogSourceId || '');
  const [existingSourceConfig, existingSourceConfigLoaded, existingSourceConfigLoadError] = state;
  const [isExpectedFormatDrawerOpen, setIsExpectedFormatDrawerOpen] = React.useState(false);

  const panelContent = (
    <DrawerPanelContent isResizable defaultSize="50%">
      <ExpectedAgentYamlFormatDrawerPanel
        onClose={() => setIsExpectedFormatDrawerOpen(false)}
      />
    </DrawerPanelContent>
  );

  return (
    <Drawer isExpanded={isExpectedFormatDrawerOpen}>
      <DrawerContent panelContent={panelContent}>
        <DrawerContentBody>
          <ApplicationsPage
            breadcrumb={
              <Breadcrumb>
                <BreadcrumbItem>
                  <Link to={agentCatalogSettingsUrl()}>
                    {AGENT_CATALOG_SETTINGS_PAGE_TITLE}
                  </Link>
                </BreadcrumbItem>
                <BreadcrumbItem data-testid="breadcrumb-agent-source-action" isActive>
                  {breadcrumbLabel}
                </BreadcrumbItem>
              </Breadcrumb>
            }
            title={pageTitle}
            description={description}
            errorMessage={catalogSourceId ? existingSourceConfigLoadError?.message : undefined}
            empty={catalogSourceId ? !existingSourceConfig : false}
            loaded={catalogSourceId ? existingSourceConfigLoaded : true}
            provideChildrenPadding
          >
            <ManageAgentSourceForm
              existingSourceConfig={existingSourceConfig || undefined}
              isEditMode={!isAddMode}
              onToggleExpectedFormatDrawer={() =>
                setIsExpectedFormatDrawerOpen((prev) => !prev)
              }
            />
          </ApplicationsPage>
        </DrawerContentBody>
      </DrawerContent>
    </Drawer>
  );
};

export default ManageAgentSourcePage;
