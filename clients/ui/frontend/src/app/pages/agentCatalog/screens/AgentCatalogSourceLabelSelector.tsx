import * as React from 'react';
import {
  Button,
  Flex,
  Stack,
  StackItem,
  Toolbar,
  ToolbarContent,
  ToolbarGroup,
  ToolbarItem,
  ToolbarToggleGroup,
} from '@patternfly/react-core';
import { ArrowRightIcon, FilterIcon } from '@patternfly/react-icons';
import { useThemeContext } from 'mod-arch-kubeflow';
import { ThemeAwareSearchInput } from 'mod-arch-shared';
import { AgentCatalogContext } from '~/app/context/agentCatalog/AgentCatalogContext';
import { hasAgentFiltersApplied } from '~/app/pages/agentCatalog/utils/agentCatalogUtils';
import AgentCatalogActiveFilters from '~/app/pages/agentCatalog/components/AgentCatalogActiveFilters';
import AgentCatalogSourceLabelBlocks from './AgentCatalogSourceLabelBlocks';

type AgentCatalogSourceLabelSelectorProps = {
  searchTerm: string;
  onSearch: (term: string) => void;
  onClearSearch: () => void;
  onResetAllFilters: () => void;
};

const AgentCatalogSourceLabelSelector: React.FC<AgentCatalogSourceLabelSelectorProps> = ({
  searchTerm,
  onSearch,
  onClearSearch,
  onResetAllFilters,
}) => {
  const [inputValue, setInputValue] = React.useState(searchTerm || '');
  const { isMUITheme } = useThemeContext();
  const { filters } = React.useContext(AgentCatalogContext);

  const hasFiltersAppliedValue = hasAgentFiltersApplied(filters, searchTerm);

  React.useEffect(() => {
    setInputValue(searchTerm || '');
  }, [searchTerm]);

  const handleClearAllFilters = React.useCallback(() => {
    if (hasFiltersAppliedValue) {
      onResetAllFilters();
    }
  }, [hasFiltersAppliedValue, onResetAllFilters]);

  const handleSearch = React.useCallback(() => {
    if (inputValue.trim() !== searchTerm) {
      onSearch(inputValue.trim());
    }
  }, [inputValue, searchTerm, onSearch]);

  const handleClear = React.useCallback(() => {
    onClearSearch();
  }, [onClearSearch]);

  const handleSearchInputChange = React.useCallback((value: string) => {
    setInputValue(value);
  }, []);

  const handleSearchInputSearch = React.useCallback(
    (_: React.SyntheticEvent<HTMLButtonElement>, value: string) => {
      onSearch(value.trim());
    },
    [onSearch],
  );

  const toolbarClearAllProps = hasFiltersAppliedValue
    ? {
        clearAllFilters: handleClearAllFilters,
        clearFiltersButtonText: 'Reset all filters' as const,
      }
    : undefined;

  return (
    <Stack hasGutter>
      <StackItem>
        <Toolbar
          className="pf-v6-u-pb-0"
          key={hasFiltersAppliedValue ? 'has-filters' : 'no-filters'}
          {...(toolbarClearAllProps ?? {})}
        >
          <ToolbarContent rowWrap={{ default: 'wrap' }}>
            <Flex style={{ flex: 1 }}>
              <ToolbarToggleGroup style={{ flex: 1 }} breakpoint="md" toggleIcon={<FilterIcon />}>
                <ToolbarGroup
                  style={{ flex: 1 }}
                  variant="filter-group"
                  gap={{ default: 'gapMd' }}
                  alignItems="center"
                >
                  <ToolbarItem style={{ flex: 1 }}>
                    <ThemeAwareSearchInput
                      data-testid="agent-catalog-search-input"
                      aria-label="Search with submit button"
                      className="toolbar-fieldset-wrapper"
                      placeholder="Search by name, keyword, or description"
                      value={inputValue}
                      onChange={handleSearchInputChange}
                      onSearch={handleSearchInputSearch}
                      onClear={handleClear}
                    />
                  </ToolbarItem>
                  <ToolbarItem>
                    {isMUITheme && (
                      <Button
                        isInline
                        aria-label="arrow-right-button"
                        data-testid="agent-search-button"
                        variant="link"
                        icon={<ArrowRightIcon />}
                        iconPosition="right"
                        onClick={handleSearch}
                      />
                    )}
                  </ToolbarItem>
                </ToolbarGroup>
              </ToolbarToggleGroup>
              {hasFiltersAppliedValue && <AgentCatalogActiveFilters />}
            </Flex>
          </ToolbarContent>
        </Toolbar>
      </StackItem>
      <StackItem>
        <Flex
          justifyContent={{ default: 'justifyContentSpaceBetween' }}
          alignItems={{ default: 'alignItemsCenter' }}
        >
          <AgentCatalogSourceLabelBlocks />
        </Flex>
      </StackItem>
    </Stack>
  );
};

export default AgentCatalogSourceLabelSelector;
