import * as React from 'react';
import { ActionsColumn, Td, Tr } from '@patternfly/react-table';
import { Button, Switch } from '@patternfly/react-core';
import { useNavigate } from 'react-router-dom';
import { CatalogSourceConfig } from '~/app/modelCatalogTypes';
import { manageAgentSourceUrl } from '~/app/routes/agentCatalogSettings/agentCatalogSettings';
import { CATALOG_SOURCE_TYPE_LABELS } from '~/concepts/modelCatalogSettings/const';
import DeleteModal from '~/app/shared/components/DeleteModal';
import { useNotification } from '~/app/hooks/useNotification';
import CatalogSourceStatus from '~/app/pages/agentCatalogSettings/components/AgentCatalogSourceStatus';

type AgentCatalogSourceConfigsTableRowProps = {
  catalogSourceConfig: CatalogSourceConfig;
  onDeleteSource: (sourceId: string) => Promise<void>;
  isUpdatingToggle: boolean;
  onToggleUpdate: (checked: boolean, sourceConfig: CatalogSourceConfig) => void;
};

const AgentCatalogSourceConfigsTableRow: React.FC<AgentCatalogSourceConfigsTableRowProps> = ({
  catalogSourceConfig,
  onDeleteSource,
  isUpdatingToggle,
  onToggleUpdate,
}) => {
  const navigate = useNavigate();
  const notification = useNotification();
  const [isDeleteModalOpen, setIsDeleteModalOpen] = React.useState(false);
  const [isDeleting, setIsDeleting] = React.useState(false);
  const [deleteError, setDeleteError] = React.useState<Error | undefined>();

  const isDefault = catalogSourceConfig.isDefault ?? false;
  const isEnabled = catalogSourceConfig.enabled ?? true;

  const handleEnableToggle = (checked: boolean) => {
    onToggleUpdate(checked, catalogSourceConfig);
  };

  const handleManageSource = () => {
    navigate(manageAgentSourceUrl(catalogSourceConfig.id));
  };

  const handleDeleteClick = () => {
    setDeleteError(undefined);
    setIsDeleteModalOpen(true);
  };

  const handleDeleteConfirm = async () => {
    setIsDeleting(true);
    setDeleteError(undefined);

    try {
      await onDeleteSource(catalogSourceConfig.id);
      setIsDeleteModalOpen(false);
      notification.success(`${catalogSourceConfig.name} deleted successfully`);
    } catch (error) {
      setDeleteError(error instanceof Error ? error : new Error('Failed to delete source'));
    } finally {
      setIsDeleting(false);
    }
  };

  const handleCloseDeleteModal = () => {
    if (!isDeleting) {
      setIsDeleteModalOpen(false);
    }
  };

  return (
    <>
      <Tr>
        <Td dataLabel="Name" style={{ verticalAlign: 'middle' }}>
          <span data-testid={`agent-source-name-${catalogSourceConfig.id}`}>
            {catalogSourceConfig.name}
          </span>
        </Td>
        <Td dataLabel="Source type" style={{ verticalAlign: 'middle' }}>
          <span data-testid={`agent-source-type-${catalogSourceConfig.id}`}>
            {CATALOG_SOURCE_TYPE_LABELS[catalogSourceConfig.type]}
          </span>
        </Td>
        <Td dataLabel="Enable" style={{ verticalAlign: 'middle' }}>
          <Switch
            data-testid={`agent-enable-toggle-${catalogSourceConfig.id}`}
            id={`agent-enable-toggle-${catalogSourceConfig.id}`}
            aria-label={`Enable ${catalogSourceConfig.name}`}
            isChecked={isEnabled}
            isDisabled={isUpdatingToggle}
            onChange={(_event, checked) => handleEnableToggle(checked)}
          />
        </Td>
        <Td dataLabel="Status" style={{ verticalAlign: 'middle' }}>
          <CatalogSourceStatus catalogSourceConfig={catalogSourceConfig} />
        </Td>
        <Td dataLabel="Actions" style={{ verticalAlign: 'middle' }}>
          <Button
            variant="link"
            onClick={handleManageSource}
            data-testid={`manage-agent-source-button-${catalogSourceConfig.id}`}
          >
            Manage source
          </Button>
        </Td>
        <Td isActionCell style={{ verticalAlign: 'middle' }}>
          {!isDefault && (
            <ActionsColumn
              items={[
                {
                  title: 'Delete source',
                  onClick: handleDeleteClick,
                },
              ]}
              data-testid={`agent-source-actions-${catalogSourceConfig.id}`}
            />
          )}
        </Td>
      </Tr>
      {isDeleteModalOpen && (
        <DeleteModal
          title="Delete a source"
          testId="delete-agent-source-modal"
          onClose={handleCloseDeleteModal}
          deleting={isDeleting}
          onDelete={handleDeleteConfirm}
          deleteName={catalogSourceConfig.name}
          error={deleteError}
        >
          The <strong>{catalogSourceConfig.name}</strong> source will be deleted, and its agents will
          be removed from the agent catalog.
        </DeleteModal>
      )}
    </>
  );
};

export default AgentCatalogSourceConfigsTableRow;
