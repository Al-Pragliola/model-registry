import * as React from 'react';
import {
  EmptyState,
  EmptyStateVariant,
  EmptyStateBody,
  EmptyStateFooter,
  EmptyStateActions,
  Flex,
  FlexItem,
  Title,
  Tabs,
  Tab,
  TabTitleText,
  Alert,
  List,
  ListItem,
  Spinner,
  Button,
  AlertActionLink,
} from '@patternfly/react-core';
import { CheckCircleIcon, TimesCircleIcon } from '@patternfly/react-icons';
import {
  UseAgentSourcePreviewResult,
  PreviewTab,
  PreviewMode,
} from '~/app/pages/agentCatalogSettings/useAgentSourcePreview';

type AgentPreviewPanelProps = {
  preview: UseAgentSourcePreviewResult;
};

const AgentPreviewPanel: React.FC<AgentPreviewPanelProps> = ({ preview }) => {
  const { previewState, handlePreview, handleTabChange, handleLoadMore, hasFormChanged, canPreview } =
    preview;
  const { isLoadingInitial, isLoadingMore, activeTab, summary, tabStates, error, mode } =
    previewState;
  const { items, hasMore } = tabStates[activeTab];
  const previewError = mode === PreviewMode.PREVIEW ? error : undefined;

  const onPreview = () => handlePreview();
  const onLoadMore = () => handleLoadMore();

  const handleTabSelect = (_event: React.MouseEvent, tabIndex: string | number) => {
    handleTabChange(tabIndex === 0 ? PreviewTab.INCLUDED : PreviewTab.EXCLUDED);
  };

  const renderEmptyState = () => {
    if (previewError) {
      return (
        <EmptyState
          icon={TimesCircleIcon}
          titleText="Preview failed"
          variant={EmptyStateVariant.sm}
        >
          <EmptyStateBody>{previewError.message}</EmptyStateBody>
          <EmptyStateFooter>
            <EmptyStateActions>
              <Button
                variant="link"
                onClick={onPreview}
                isDisabled={!canPreview}
                isLoading={isLoadingInitial}
                data-testid="agent-preview-button-retry"
              >
                Preview
              </Button>
            </EmptyStateActions>
          </EmptyStateFooter>
        </EmptyState>
      );
    }

    return (
      <EmptyState titleText="Preview agents" variant={EmptyStateVariant.sm}>
        <EmptyStateBody>
          To view the agents from this source that will appear in the agent catalog, complete all
          required fields, then click <strong>Preview</strong>.
        </EmptyStateBody>
        <EmptyStateFooter>
          <EmptyStateActions>
            <Button
              variant="link"
              onClick={onPreview}
              isDisabled={!canPreview}
              isLoading={isLoadingInitial}
              data-testid="agent-preview-button-panel"
            >
              Preview
            </Button>
          </EmptyStateActions>
        </EmptyStateFooter>
      </EmptyState>
    );
  };

  const renderContent = () => {
    if (isLoadingInitial) {
      return (
        <div className="pf-v6-u-text-align-center pf-v6-u-py-xl">
          <Spinner size="xl" aria-label="Loading preview" />
        </div>
      );
    }

    if ((!items.length && !summary) || previewError) {
      return renderEmptyState();
    }

    return (
      <>
        <Tabs
          activeKey={activeTab === PreviewTab.INCLUDED ? 0 : 1}
          onSelect={handleTabSelect}
          aria-label="Preview tabs"
        >
          <Tab eventKey={0} title={<TabTitleText>Agents included</TabTitleText>} />
          <Tab eventKey={1} title={<TabTitleText>Agents excluded</TabTitleText>} />
        </Tabs>
        <div className="pf-v6-u-mt-md">
          {hasFormChanged && (
            <Alert
              variant="info"
              isInline
              title="Source configuration changed. Refresh the preview."
              className="pf-v6-u-mb-md"
              actionLinks={
                <AlertActionLink onClick={onPreview} data-testid="agent-refresh-preview-link">
                  Refresh preview
                </AlertActionLink>
              }
            />
          )}
          {items.length > 0 ? (
            <>
              <strong>
                {activeTab === PreviewTab.INCLUDED
                  ? `${summary?.includedModels ?? 0} of ${summary?.totalModels ?? 0} agents included:`
                  : `${summary?.excludedModels ?? 0} of ${summary?.totalModels ?? 0} agents excluded:`}
              </strong>
              <List isPlain className="pf-v6-u-mt-md">
                {items.map((agent) => (
                  <ListItem
                    key={agent.name}
                    icon={
                      agent.included ? (
                        <CheckCircleIcon color="green" />
                      ) : (
                        <TimesCircleIcon color="red" />
                      )
                    }
                  >
                    {agent.name}
                  </ListItem>
                ))}
              </List>
              {hasMore && (
                <div className="pf-v6-u-mt-md pf-v6-u-text-align-center">
                  <Button
                    variant="link"
                    onClick={onLoadMore}
                    isLoading={isLoadingMore}
                    isDisabled={isLoadingMore}
                  >
                    {isLoadingMore ? 'Loading...' : 'Load more'}
                  </Button>
                </div>
              )}
            </>
          ) : (
            <EmptyState
              variant={EmptyStateVariant.sm}
              titleText={
                activeTab === PreviewTab.INCLUDED ? 'No agents included' : 'No agents excluded'
              }
            >
              <EmptyStateBody>
                {activeTab === PreviewTab.INCLUDED
                  ? 'No agents from this source match the inclusion filters.'
                  : 'No agents from this source are excluded by the filters.'}
              </EmptyStateBody>
            </EmptyState>
          )}
        </div>
      </>
    );
  };

  return (
    <div data-testid="agent-preview-panel" className="pf-v6-u-h-100">
      <Flex
        justifyContent={{ default: 'justifyContentSpaceBetween' }}
        alignItems={{ default: 'alignItemsCenter' }}
        className="pf-v6-u-mb-md"
      >
        <FlexItem>
          <Title headingLevel="h2" size="lg">
            Agent catalog preview
          </Title>
        </FlexItem>
        <FlexItem>
          <Button
            variant="secondary"
            onClick={onPreview}
            isDisabled={!canPreview}
            isLoading={isLoadingInitial}
            data-testid="agent-preview-button-header"
          >
            Preview
          </Button>
        </FlexItem>
      </Flex>
      {renderContent()}
    </div>
  );
};

export default AgentPreviewPanel;
