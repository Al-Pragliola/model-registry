import React from 'react';
import { AgentCatalogContext } from '~/app/context/agentCatalog/AgentCatalogContext';
import { CatalogAllItemsView } from '~/app/shared/components/catalog';
import AgentCatalogCategorySection from './AgentCatalogCategorySection';

type AgentCatalogAllAgentsViewProps = {
  searchTerm: string;
};

const AgentCatalogAllAgentsView: React.FC<AgentCatalogAllAgentsViewProps> = ({ searchTerm }) => {
  const { catalogSources, catalogLabels, setSelectedSourceLabel } =
    React.useContext(AgentCatalogContext);

  const handleShowMoreCategory = React.useCallback(
    (categoryLabel: string) => {
      setSelectedSourceLabel(categoryLabel);
    },
    [setSelectedSourceLabel],
  );

  return (
    <CatalogAllItemsView
      searchTerm={searchTerm}
      catalogSources={catalogSources}
      catalogLabels={catalogLabels}
      pageSize={4}
      otherSectionKey="other-agents"
      onShowMore={handleShowMoreCategory}
      renderCategorySection={(label, term, pageSize, onShowMore) => (
        <AgentCatalogCategorySection
          label={label}
          searchTerm={term}
          pageSize={pageSize}
          onShowMore={onShowMore}
        />
      )}
    />
  );
};

export default AgentCatalogAllAgentsView;
