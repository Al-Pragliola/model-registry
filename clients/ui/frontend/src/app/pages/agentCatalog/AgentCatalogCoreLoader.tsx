import * as React from 'react';
import { Alert, Bullseye } from '@patternfly/react-core';
import {
  ApplicationsPage,
  KubeflowDocs,
  ProjectObjectType,
  TitleWithIcon,
  typedEmptyImage,
  WhosMyAdministrator,
} from 'mod-arch-shared';
import { useThemeContext } from 'mod-arch-kubeflow';
import { Outlet } from 'react-router-dom';
import { AgentCatalogContext } from '~/app/context/agentCatalog/AgentCatalogContext';
import { EmptyCatalogState, hasSourcesWithModels } from '~/app/shared/components/catalog';
import { AGENT_CATALOG_TITLE, AGENT_CATALOG_DESCRIPTION } from '~/app/pages/agentCatalog/const';

const AgentCatalogCoreLoader: React.FC = () => {
  const { catalogSources, catalogSourcesLoaded, catalogSourcesLoadError } =
    React.useContext(AgentCatalogContext);
  const { isMUITheme } = useThemeContext();

  if (catalogSourcesLoadError) {
    return (
      <ApplicationsPage
        title={
          <TitleWithIcon
            title={AGENT_CATALOG_TITLE}
            objectType={ProjectObjectType.mcpCatalog}
          />
        }
        description={AGENT_CATALOG_DESCRIPTION}
        headerContent={null}
        empty
        emptyStatePage={
          <Bullseye>
            <Alert title="Agent catalog source load error" variant="danger" isInline>
              {catalogSourcesLoadError.message}
            </Alert>
          </Bullseye>
        }
        loaded
      />
    );
  }

  if (!catalogSourcesLoaded) {
    return (
      <ApplicationsPage
        title={
          <TitleWithIcon
            title={AGENT_CATALOG_TITLE}
            objectType={ProjectObjectType.mcpCatalog}
          />
        }
        description={AGENT_CATALOG_DESCRIPTION}
        headerContent={null}
        empty
        emptyStatePage={<Bullseye>Loading catalog sources...</Bullseye>}
        loaded={false}
      />
    );
  }

  if (catalogSources?.items?.length === 0 || !hasSourcesWithModels(catalogSources)) {
    return (
      <ApplicationsPage
        title={
          <TitleWithIcon
            title={AGENT_CATALOG_TITLE}
            objectType={ProjectObjectType.mcpCatalog}
          />
        }
        description={AGENT_CATALOG_DESCRIPTION}
        empty
        emptyStatePage={
          <EmptyCatalogState
            testid="empty-agent-catalog-state"
            title="Agent catalog configuration required"
            description={
              isMUITheme
                ? 'To discover agents, follow the instructions in the docs below.'
                : 'There are no agent sources to display. Request that your administrator configure agent sources for the catalog.'
            }
            headerIcon={() => (
              <img src={typedEmptyImage(ProjectObjectType.modelRegistrySettings)} alt="" />
            )}
            primaryAction={isMUITheme ? <KubeflowDocs /> : <WhosMyAdministrator />}
          />
        }
        headerContent={null}
        loaded
        provideChildrenPadding
      />
    );
  }

  return <Outlet />;
};

export default AgentCatalogCoreLoader;
