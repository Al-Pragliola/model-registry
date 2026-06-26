import * as React from 'react';
import { PageSection, Stack, StackItem } from '@patternfly/react-core';
import { ApplicationsPage, ProjectObjectType, TitleWithIcon } from 'mod-arch-shared';
import { SearchIcon } from '@patternfly/react-icons';
import { AgentCatalogContext } from '~/app/context/agentCatalog/AgentCatalogContext';
import { hasAgentFiltersApplied } from '~/app/pages/agentCatalog/utils/agentCatalogUtils';
import { AGENT_CATALOG_TITLE, AGENT_CATALOG_DESCRIPTION } from '~/app/pages/agentCatalog/const';
import { EmptyCatalogState } from '~/app/shared/components/catalog';
import { getActiveSourceLabels } from '~/app/shared/components/catalog/utils/catalogSourceUtils';
import ScrollViewOnMount from '~/app/shared/components/ScrollViewOnMount';
import AgentCatalogSourceLabelSelector from './AgentCatalogSourceLabelSelector';
import AgentCatalogAllAgentsView from './AgentCatalogAllAgentsView';
import AgentCatalogGalleryView from './AgentCatalogGalleryView';

const AgentCatalog: React.FC = () => {
  const {
    searchQuery,
    setSearchQuery,
    clearAllFilters,
    selectedSourceLabel,
    setSelectedSourceLabel,
    filters,
    catalogSources,
    catalogLabels,
    catalogSourcesLoaded,
  } = React.useContext(AgentCatalogContext);

  const filtersApplied = hasAgentFiltersApplied(filters, searchQuery);
  const isAllAgentsView = selectedSourceLabel === undefined && !filtersApplied;

  const activeCategories = React.useMemo(
    () => getActiveSourceLabels(catalogSources, catalogLabels),
    [catalogSources, catalogLabels],
  );

  const isSingleCategory = activeCategories.length === 1;
  const hasNoCategories = activeCategories.length === 0;

  React.useEffect(() => {
    if (catalogSourcesLoaded && isSingleCategory && selectedSourceLabel !== activeCategories[0]) {
      setSelectedSourceLabel(activeCategories[0]);
    }
  }, [
    catalogSourcesLoaded,
    isSingleCategory,
    activeCategories,
    selectedSourceLabel,
    setSelectedSourceLabel,
  ]);

  const handleSearch = React.useCallback(
    (term: string) => {
      setSearchQuery(term);
    },
    [setSearchQuery],
  );

  const handleClearSearch = React.useCallback(() => {
    setSearchQuery('');
  }, [setSearchQuery]);

  const handleResetAllFilters = React.useCallback(() => {
    clearAllFilters();
  }, [clearAllFilters]);

  return (
    <ApplicationsPage
      title={
        <TitleWithIcon title={AGENT_CATALOG_TITLE} objectType={ProjectObjectType.mcpCatalog} />
      }
      description={AGENT_CATALOG_DESCRIPTION}
      empty={false}
      loaded
      provideChildrenPadding
    >
      <ScrollViewOnMount shouldScroll scrollToTop />
      {catalogSourcesLoaded && hasNoCategories ? (
        <EmptyCatalogState
          testid="empty-agent-catalog-no-categories"
          title="No agents available"
          headerIcon={SearchIcon}
          description="There are no agent categories available. Configure sources in settings to get started."
        />
      ) : (
        <Stack hasGutter>
          <StackItem>
            <AgentCatalogSourceLabelSelector
              searchTerm={searchQuery}
              onSearch={handleSearch}
              onClearSearch={handleClearSearch}
              onResetAllFilters={handleResetAllFilters}
            />
          </StackItem>
          <StackItem isFilled>
            {isAllAgentsView && !isSingleCategory ? (
              <AgentCatalogAllAgentsView searchTerm={searchQuery} />
            ) : (
              <AgentCatalogGalleryView
                handleFilterReset={handleResetAllFilters}
                isSingleCategory={isSingleCategory}
                singleCategoryLabel={isSingleCategory ? activeCategories[0] : undefined}
              />
            )}
          </StackItem>
        </Stack>
      )}
    </ApplicationsPage>
  );
};

export default AgentCatalog;
