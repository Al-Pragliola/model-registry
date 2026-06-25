import * as React from 'react';
import { Button, Label, Spinner, Stack, StackItem, Truncate } from '@patternfly/react-core';
import { InProgressIcon } from '@patternfly/react-icons';
import { CatalogSourceConfig } from '~/app/modelCatalogTypes';
import { AgentCatalogSettingsContext } from '~/app/context/agentCatalogSettings/AgentCatalogSettingsContext';
import { CatalogSourceStatus as CatalogSourceStatusEnum } from '~/concepts/modelCatalogSettings/const';
import CatalogSourceStatusErrorModal from '~/app/pages/modelCatalogSettings/components/CatalogSourceStatusErrorModal';

type AgentCatalogSourceStatusProps = {
  catalogSourceConfig: CatalogSourceConfig;
};

const AgentCatalogSourceStatus: React.FC<AgentCatalogSourceStatusProps> = ({
  catalogSourceConfig,
}) => {
  const { catalogSources, catalogSourcesLoaded, catalogSourcesLoadError } = React.useContext(
    AgentCatalogSettingsContext,
  );
  const [isErrorModalOpen, setIsErrorModalOpen] = React.useState(false);

  // Don't render status for default sources
  if (catalogSourceConfig.isDefault) {
    return <>-</>;
  }
  // If source is disabled, render "-"
  if (!catalogSourceConfig.enabled) {
    return <>-</>;
  }

  // Show loading spinner while fetching sources
  if (!catalogSourcesLoaded) {
    return (
      <Spinner size="md" data-testid={`agent-source-status-loading-${catalogSourceConfig.id}`} />
    );
  }

  // Find the matching source from the catalog sources list
  const matchingSource = catalogSources?.items?.find(
    (source) => source.id === catalogSourceConfig.id,
  );

  const startingOrUnknownLabel = (
    <Label
      color="grey"
      variant="outline"
      icon={<InProgressIcon />}
      data-testid={`agent-source-status-${catalogSourcesLoadError ? 'unknown' : 'starting'}-${catalogSourceConfig.id}`}
    >
      {catalogSourcesLoadError ? 'Unknown' : 'Starting'}
    </Label>
  );

  if (!matchingSource || !matchingSource.status) {
    return startingOrUnknownLabel;
  }

  // Render based on status
  switch (matchingSource.status) {
    case CatalogSourceStatusEnum.AVAILABLE:
      return (
        <Label
          status="success"
          variant="outline"
          data-testid={`agent-source-status-connected-${catalogSourceConfig.id}`}
        >
          Ready
        </Label>
      );

    case CatalogSourceStatusEnum.ERROR: {
      const errorMessage = matchingSource.error || 'Unknown error occurred';

      return (
        <>
          <Stack hasGutter>
            <StackItem>
              <Label
                status="danger"
                variant="outline"
                data-testid={`agent-source-status-failed-${catalogSourceConfig.id}`}
              >
                Failed
              </Label>
            </StackItem>
            <StackItem>
              <Button
                variant="link"
                isInline
                isDanger
                onClick={() => setIsErrorModalOpen(true)}
                data-testid={`agent-source-status-error-link-${catalogSourceConfig.id}`}
              >
                <Truncate content={errorMessage} tooltipProps={{ hidden: true }} />
              </Button>
            </StackItem>
          </Stack>
          <CatalogSourceStatusErrorModal
            isOpen={isErrorModalOpen}
            onClose={() => setIsErrorModalOpen(false)}
            errorMessage={errorMessage}
          />
        </>
      );
    }

    case CatalogSourceStatusEnum.DISABLED:
      // If we reach here, config.enabled is true
      // But status is still DISABLED, so show "Starting" (re-enable case)
      return startingOrUnknownLabel;
    default:
      return startingOrUnknownLabel;
  }
};

export default AgentCatalogSourceStatus;
