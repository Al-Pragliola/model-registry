import * as React from 'react';
import { AgentCatalogContext } from '~/app/context/agentCatalog/AgentCatalogContext';
import { CatalogFilterPanel, useCatalogFilterConfigs } from '~/app/shared/components/catalog';
import { AGENT_FILTER_KEYS, AGENT_FILTER_CATEGORY_NAMES } from '~/app/pages/agentCatalog/const';

const AgentCatalogFilters: React.FC = () => {
  const { filters, setFilters, filterOptions, filterOptionsLoaded, filterOptionsLoadError } =
    React.useContext(AgentCatalogContext);

  const onFilterChange = React.useCallback(
    (key: string, values: string[]) => {
      setFilters((prev) => ({ ...prev, [key]: values }));
    },
    [setFilters],
  );

  const filterPanelItems = useCatalogFilterConfigs({
    filterKeys: AGENT_FILTER_KEYS,
    filterNames: AGENT_FILTER_CATEGORY_NAMES,
    filterOptions: filterOptions?.filters,
    selectedFilters: filters,
    onFilterChange,
  });

  return (
    <CatalogFilterPanel
      loaded={filterOptionsLoaded}
      loadError={filterOptionsLoadError}
      filters={filterPanelItems}
      testIdPrefix="agent-filter"
    />
  );
};

export default AgentCatalogFilters;
