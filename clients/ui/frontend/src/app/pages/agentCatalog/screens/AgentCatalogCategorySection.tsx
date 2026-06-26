import * as React from 'react';
import { useAgentsBySourceLabelWithAPI } from '~/app/hooks/agentCatalog/useAgentsBySourceLabel';
import {
  getLabelDescription,
  getLabelDisplayName,
  CatalogCategorySection,
} from '~/app/shared/components/catalog';
import { AgentCatalogContext } from '~/app/context/agentCatalog/AgentCatalogContext';
import {
  AGENT_CATALOG_GRID_SPAN,
  OTHER_AGENTS_DISPLAY_NAME,
} from '~/app/pages/agentCatalog/const';
import AgentCatalogCard from '~/app/pages/agentCatalog/components/AgentCatalogCard';

type AgentCatalogCategorySectionProps = {
  label: string;
  searchTerm: string;
  pageSize: number;
  onShowMore: (label: string) => void;
};

const AgentCatalogCategorySection: React.FC<AgentCatalogCategorySectionProps> = ({
  label,
  searchTerm,
  pageSize,
  onShowMore,
}) => {
  const { agentApiState, catalogLabels } = React.useContext(AgentCatalogContext);
  const { agents, agentsLoaded, agentsLoadError } = useAgentsBySourceLabelWithAPI(agentApiState, {
    sourceLabel: label,
    pageSize,
    searchQuery: searchTerm,
  });

  const categoryTitle = getLabelDisplayName(
    label,
    catalogLabels,
    OTHER_AGENTS_DISPLAY_NAME,
    'agents',
  );
  const categoryDescription = getLabelDescription(label, catalogLabels);
  const labelSlug = label.toLowerCase().replace(/\s+/g, '-');

  return (
    <CatalogCategorySection
      label={label}
      categoryTitle={categoryTitle}
      categoryDescription={categoryDescription}
      items={agents.items}
      loaded={agentsLoaded}
      loadError={agentsLoadError}
      pageSize={pageSize}
      onShowMore={onShowMore}
      renderCard={(agent) => <AgentCatalogCard agent={agent} />}
      getItemKey={(agent) => agent.id}
      gridSpans={AGENT_CATALOG_GRID_SPAN}
      loadingScreenReaderText={`Loading ${label} agents`}
      testIds={{
        title: `agent-category-title-${label}`,
        showMore: `agent-show-all-${labelSlug}`,
        error: `agent-error-state-${label}`,
        skeleton: (index) => `agent-category-skeleton-${labelSlug}-${index}`,
        empty: `empty-agent-catalog-state-${label}`,
      }}
    />
  );
};
export default AgentCatalogCategorySection;
