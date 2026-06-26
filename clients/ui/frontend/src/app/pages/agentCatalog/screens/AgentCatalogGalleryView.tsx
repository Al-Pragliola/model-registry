import * as React from 'react';
import { Button } from '@patternfly/react-core';
import { SearchIcon } from '@patternfly/react-icons';
import { AgentCatalogContext } from '~/app/context/agentCatalog/AgentCatalogContext';
import { useAgentsBySourceLabelWithAPI } from '~/app/hooks/agentCatalog/useAgentsBySourceLabel';
import {
  AGENT_CATALOG_GRID_SPAN,
  OTHER_AGENTS_DISPLAY_NAME,
} from '~/app/pages/agentCatalog/const';
import { agentFiltersToFilterQuery } from '~/app/pages/agentCatalog/utils/agentCatalogUtils';
import {
  getLabelDisplayName,
  getLabelDescription,
  CatalogGalleryLayout,
  EmptyCatalogState,
} from '~/app/shared/components/catalog';
import AgentCatalogCard from '~/app/pages/agentCatalog/components/AgentCatalogCard';

const PAGE_SIZE = 10;

type AgentCatalogGalleryViewProps = {
  handleFilterReset: () => void;
  isSingleCategory?: boolean;
  singleCategoryLabel?: string;
};

const AgentCatalogGalleryView: React.FC<AgentCatalogGalleryViewProps> = ({
  handleFilterReset,
  isSingleCategory = false,
  singleCategoryLabel,
}) => {
  const {
    agentApiState,
    selectedSourceLabel,
    searchQuery,
    filters,
    catalogLabels,
    catalogLabelsLoaded,
  } = React.useContext(AgentCatalogContext);

  const filterQuery = React.useMemo(() => agentFiltersToFilterQuery(filters), [filters]);

  const { agents, agentsLoaded, agentsLoadError } = useAgentsBySourceLabelWithAPI(agentApiState, {
    sourceLabel: selectedSourceLabel,
    pageSize: PAGE_SIZE,
    searchQuery,
    filterQuery: filterQuery || undefined,
  });

  const loaded = agentsLoaded && catalogLabelsLoaded;

  const effectiveCategoryLabel = singleCategoryLabel || selectedSourceLabel || '';
  const categoryTitle = isSingleCategory
    ? getLabelDisplayName(
        effectiveCategoryLabel,
        catalogLabels,
        OTHER_AGENTS_DISPLAY_NAME,
        'agents',
      )
    : undefined;
  const categoryDescription = isSingleCategory
    ? getLabelDescription(effectiveCategoryLabel, catalogLabels)
    : undefined;

  return (
    <CatalogGalleryLayout
      items={agents.items}
      loaded={loaded}
      loadError={agentsLoadError}
      renderCard={(agent) => <AgentCatalogCard agent={agent} />}
      getItemKey={(agent) => agent.id}
      gridSpans={AGENT_CATALOG_GRID_SPAN}
      hasMore={agents.hasMore && agents.items.length >= PAGE_SIZE}
      isLoadingMore={agents.isLoadingMore}
      onLoadMore={agents.loadMore}
      loadMoreLabel="Load more agents"
      loadingMoreLabel="Loading more agents..."
      loadingLabel="Loading agents..."
      errorTitle="Failed to load agents"
      categoryTitle={categoryTitle}
      categoryDescription={categoryDescription}
      renderEmptyState={() => (
        <EmptyCatalogState
          testid="empty-agent-catalog-state"
          title="No results found"
          headerIcon={SearchIcon}
          description="Adjust your filters and try again."
          primaryAction={
            <Button variant="link" onClick={handleFilterReset}>
              Reset filters
            </Button>
          }
        />
      )}
    />
  );
};

export default AgentCatalogGalleryView;
