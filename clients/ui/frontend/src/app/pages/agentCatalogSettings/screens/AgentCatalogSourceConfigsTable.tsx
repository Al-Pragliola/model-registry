import * as React from 'react';
import {
  Alert,
  AlertActionCloseButton,
  Button,
  Flex,
  FlexItem,
  Stack,
  StackItem,
  Toolbar,
  ToolbarContent,
  ToolbarItem,
} from '@patternfly/react-core';
import { Table } from 'mod-arch-shared';
import { CatalogSourceConfig } from '~/app/modelCatalogTypes';
import { AgentCatalogSettingsContext } from '~/app/context/agentCatalogSettings/AgentCatalogSettingsContext';
import { ADD_AGENT_SOURCE_TITLE } from '~/app/routes/agentCatalogSettings/agentCatalogSettings';
import { agentCatalogSourceConfigsColumns } from './AgentCatalogSourceConfigsTableColumns';
import AgentCatalogSourceConfigsTableRow from './AgentCatalogSourceConfigsTableRow';

type AgentCatalogSourceConfigsTableProps = {
  catalogSourceConfigs: CatalogSourceConfig[];
  onAddSource: () => void;
  onDeleteSource: (sourceId: string) => Promise<void>;
};

const AgentCatalogSourceConfigsTable: React.FC<AgentCatalogSourceConfigsTableProps> = ({
  catalogSourceConfigs,
  onAddSource,
  onDeleteSource,
}) => {
  const [toggleError, setToggleError] = React.useState<Error | undefined>(undefined);
  const [updatingToggleId, setUpdatingToggleId] = React.useState<string | null>(null);
  const { apiState, refreshCatalogSourceConfigs, catalogSourcesLoadError } = React.useContext(
    AgentCatalogSettingsContext,
  );

  const handleEnableToggle = async (checked: boolean, catalogSourceConfig: CatalogSourceConfig) => {
    if (!apiState.apiAvailable) {
      setToggleError(new Error('API is not available'));
      return;
    }
    setUpdatingToggleId(catalogSourceConfig.id);
    setToggleError(undefined);

    try {
      await apiState.api.updateCatalogSourceConfig({}, catalogSourceConfig.id, {
        enabled: checked,
      });
      setToggleError(undefined);
      refreshCatalogSourceConfigs();
    } catch (e) {
      if (e instanceof Error) {
        setToggleError(new Error(`Error enabling/disabling source ${catalogSourceConfig.name}`));
      }
    } finally {
      setUpdatingToggleId(null);
    }
  };

  return (
    <Stack hasGutter>
      {catalogSourcesLoadError && (
        <StackItem>
          <Alert
            variant="danger"
            isInline
            title="Error fetching source statuses"
            data-testid="agent-source-status-error-alert"
          >
            {catalogSourcesLoadError.message}
          </Alert>
        </StackItem>
      )}
      <StackItem>
        <Table
          data-testid="agent-catalog-source-configs-table"
          data={catalogSourceConfigs}
          columns={agentCatalogSourceConfigsColumns}
          toolbarContent={
            <Flex direction={{ default: 'column' }}>
              <FlexItem>
                <Toolbar>
                  <ToolbarContent>
                    <ToolbarItem>
                      <Button
                        variant="primary"
                        onClick={onAddSource}
                        data-testid="add-agent-source-button"
                      >
                        {ADD_AGENT_SOURCE_TITLE}
                      </Button>
                    </ToolbarItem>
                  </ToolbarContent>
                </Toolbar>
              </FlexItem>
              {toggleError && (
                <FlexItem>
                  <Alert
                    variant="danger"
                    data-testid="agent-toggle-alert"
                    title={toggleError.message}
                    actionClose={
                      <AlertActionCloseButton onClose={() => setToggleError(undefined)} />
                    }
                  />
                </FlexItem>
              )}
            </Flex>
          }
          rowRenderer={(config) => (
            <AgentCatalogSourceConfigsTableRow
              key={config.id}
              catalogSourceConfig={config}
              isUpdatingToggle={updatingToggleId === config.id}
              onToggleUpdate={handleEnableToggle}
              onDeleteSource={onDeleteSource}
            />
          )}
          variant="compact"
        />
      </StackItem>
    </Stack>
  );
};

export default AgentCatalogSourceConfigsTable;
