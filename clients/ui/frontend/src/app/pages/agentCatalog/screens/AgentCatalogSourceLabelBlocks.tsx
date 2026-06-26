import * as React from 'react';
import { AgentCatalogContext } from '~/app/context/agentCatalog/AgentCatalogContext';
import { CatalogSourceLabelToggle, getLabelDisplayName } from '~/app/shared/components/catalog';
import { OTHER_AGENTS_DISPLAY_NAME } from '~/app/pages/agentCatalog/const';

const ALL_AGENTS_LABEL = 'All agents';

const AgentCatalogSourceLabelBlocks: React.FC = () => {
  const { catalogSources, catalogLabels, selectedSourceLabel, setSelectedSourceLabel } =
    React.useContext(AgentCatalogContext);

  const getLabelDisplayNameForAgent = React.useCallback(
    (label: string) =>
      getLabelDisplayName(label, catalogLabels, OTHER_AGENTS_DISPLAY_NAME, 'agents'),
    [catalogLabels],
  );

  return (
    <CatalogSourceLabelToggle
      catalogSources={catalogSources}
      catalogLabels={catalogLabels}
      selectedSourceLabel={selectedSourceLabel}
      onSelectSourceLabel={setSelectedSourceLabel}
      allBlockLabel={undefined}
      allBlockDisplayName={ALL_AGENTS_LABEL}
      testId="agent-catalog-category-toggle"
      ariaLabel="Agent category selection"
      hideWhenSingleCategory
      getLabelDisplayNameOverride={getLabelDisplayNameForAgent}
      getTestId={(blockId) => `agent-category-${blockId}`}
    />
  );
};

export default AgentCatalogSourceLabelBlocks;
