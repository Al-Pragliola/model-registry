import React from 'react';
import { useParams } from 'react-router';
import { Link } from 'react-router-dom';
import {
  Breadcrumb,
  BreadcrumbItem,
  Button,
  Content,
  ContentVariants,
  EmptyState,
  EmptyStateBody,
  EmptyStateFooter,
  Flex,
  FlexItem,
  Label,
  Stack,
  StackItem,
} from '@patternfly/react-core';
import { AutomationIcon, SearchIcon } from '@patternfly/react-icons';
import { ApplicationsPage } from 'mod-arch-shared';
import { useAgentWithAPI } from '~/app/hooks/agentCatalog/useAgent';
import { AgentCatalogContext } from '~/app/context/agentCatalog/AgentCatalogContext';
import { agentCatalogUrl } from '~/app/routes/agentCatalog/agentCatalog';
import ScrollViewOnMount from '~/app/shared/components/ScrollViewOnMount';
import AgentDetailsView from './AgentDetailsView';

const AgentDetailsPage: React.FC = () => {
  const { agentId = '' } = useParams<{ agentId: string }>();
  const { agentApiState } = React.useContext(AgentCatalogContext);
  const [agent, agentLoaded, agentLoadError] = useAgentWithAPI(agentApiState, agentId);

  const isNotFound = !agent && (agentLoaded || !!agentLoadError);

  return (
    <>
      <ScrollViewOnMount shouldScroll scrollToTop />
      <ApplicationsPage
        breadcrumb={
          <Breadcrumb>
            <BreadcrumbItem>
              <Link to={agentCatalogUrl()}>Agent Catalog</Link>
            </BreadcrumbItem>
            <BreadcrumbItem isActive data-testid="breadcrumb-agent-name">
              {agent?.displayName || agent?.name || 'Details'}
            </BreadcrumbItem>
          </Breadcrumb>
        }
        title={
          agent ? (
            <Flex
              spaceItems={{ default: 'spaceItemsMd' }}
              alignItems={{ default: 'alignItemsCenter' }}
            >
              {agent.logo ? (
                <img
                  src={agent.logo}
                  alt="agent logo"
                  style={{ height: '56px', width: '56px' }}
                />
              ) : (
                <AutomationIcon
                  style={{ fontSize: '56px' }}
                  data-testid="agent-default-icon"
                />
              )}
              <Stack>
                <StackItem>
                  <Flex
                    gap={{ default: 'gapSm' }}
                    alignItems={{ default: 'alignItemsCenter' }}
                    flexWrap={{ default: 'wrap' }}
                  >
                    <FlexItem>{agent.displayName || agent.name}</FlexItem>
                    {agent.framework && (
                      <FlexItem>
                        <Label data-testid="agent-details-framework-label">
                          {agent.framework}
                        </Label>
                      </FlexItem>
                    )}
                  </Flex>
                </StackItem>
                {agent.agentType && (
                  <StackItem>
                    <Content component={ContentVariants.small}>Type: {agent.agentType}</Content>
                  </StackItem>
                )}
              </Stack>
            </Flex>
          ) : null
        }
        empty={isNotFound}
        emptyStatePage={
          isNotFound ? (
            <EmptyState
              icon={SearchIcon}
              titleText="Agent not found"
              data-testid="agent-not-found"
            >
              <EmptyStateBody>The requested agent could not be found.</EmptyStateBody>
              <EmptyStateFooter>
                <Button
                  variant="primary"
                  component={(props) => <Link {...props} to={agentCatalogUrl()} />}
                >
                  Return to Agent Catalog
                </Button>
              </EmptyStateFooter>
            </EmptyState>
          ) : undefined
        }
        loadError={isNotFound ? undefined : agentLoadError}
        loaded={isNotFound || agentLoaded}
        errorMessage="Unable to load agent details"
        provideChildrenPadding
      >
        {agent && <AgentDetailsView agent={agent} />}
      </ApplicationsPage>
    </>
  );
};

export default AgentDetailsPage;
