import * as React from 'react';
import { ToolbarFilter, ToolbarLabel, ToolbarLabelGroup } from '@patternfly/react-core';
import { AgentCatalogContext } from '~/app/context/agentCatalog/AgentCatalogContext';
import type { AgentFilterCategoryKey } from '~/app/pages/agentCatalog/types/agentCatalogFilterOptions';
import { AGENT_FILTER_KEYS, AGENT_FILTER_CATEGORY_NAMES } from '~/app/pages/agentCatalog/const';

const AgentCatalogActiveFilters: React.FC = () => {
  const { filters, setFilters } = React.useContext(AgentCatalogContext);

  const handleRemoveFilter = React.useCallback(
    (categoryKey: AgentFilterCategoryKey, valueKey: string) => {
      setFilters((prev) => {
        const current = prev[categoryKey];
        const arr = Array.isArray(current) ? current : [];
        const newValues = arr.filter((v) => v !== valueKey);
        return { ...prev, [categoryKey]: newValues };
      });
    },
    [setFilters],
  );

  const handleClearCategory = React.useCallback(
    (categoryKey: AgentFilterCategoryKey) => {
      setFilters((prev) => ({ ...prev, [categoryKey]: [] }));
    },
    [setFilters],
  );

  return (
    <>
      {AGENT_FILTER_KEYS.map((filterKey) => {
        const filterValue = filters[filterKey];
        const values = Array.isArray(filterValue) ? filterValue : [];
        const hasValue = values.length > 0;

        const labels: ToolbarLabel[] = hasValue
          ? values.map((value) => ({
              key: value,
              node: <span data-testid={`agent-filter-chip-${filterKey}-${value}`}>{value}</span>,
            }))
          : [];

        const categoryLabelGroup: ToolbarLabelGroup = {
          key: filterKey,
          name: AGENT_FILTER_CATEGORY_NAMES[filterKey],
        };

        return (
          <ToolbarFilter
            key={filterKey}
            categoryName={categoryLabelGroup}
            labels={labels}
            deleteLabel={(_, label) => {
              const labelKey = typeof label === 'string' ? label : label.key;
              handleRemoveFilter(filterKey, labelKey);
            }}
            deleteLabelGroup={() => handleClearCategory(filterKey)}
            data-testid={`agent-filter-container-${filterKey}`}
          >
            {null}
          </ToolbarFilter>
        );
      })}
    </>
  );
};

export default AgentCatalogActiveFilters;
