import React from 'react';
import {
  Dropdown,
  DropdownItem,
  DropdownList,
  MenuToggle,
  MenuToggleElement,
  Spinner,
} from '@patternfly/react-core';
import { FileExportIcon } from '@patternfly/react-icons';
import { useQueryParamNamespaces } from 'mod-arch-core';
import { CatalogModel, ModelCatalogFilterStates, CatalogFilterOptionsList } from '~/app/modelCatalogTypes';
import { filtersToFilterQuery } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import {
  generateCatalogModelCSV,
  downloadCSV,
  getExportFilename,
} from '~/app/pages/modelCatalog/utils/csvExportUtils';
import { URL_PREFIX, BFF_API_VERSION } from '~/app/utilities/const';

type ModelCatalogExportButtonProps = {
  selectedModels: Map<string, CatalogModel>;
  clearSelection: () => void;
  selectedSourceLabel?: string;
  searchTerm?: string;
  filterData?: ModelCatalogFilterStates;
  filterOptions?: CatalogFilterOptionsList | null;
};

const ModelCatalogExportButton: React.FC<ModelCatalogExportButtonProps> = ({
  selectedModels,
  clearSelection,
  selectedSourceLabel,
  searchTerm,
  filterData,
  filterOptions,
}) => {
  const [isOpen, setIsOpen] = React.useState(false);
  const [isExporting, setIsExporting] = React.useState(false);
  const queryParams = useQueryParamNamespaces();

  const buildExportUrl = React.useCallback(
    (includeFilters: boolean): string => {
      const params = new URLSearchParams();

      if (queryParams) {
        for (const [key, value] of Object.entries(queryParams)) {
          if (value != null) {
            params.set(key, String(value));
          }
        }
      }

      if (includeFilters) {
        if (selectedSourceLabel && selectedSourceLabel !== 'All models') {
          params.set('sourceLabel', selectedSourceLabel);
        }
        if (searchTerm) {
          params.set('q', searchTerm);
        }
        if (filterData && filterOptions) {
          const filterQuery = filtersToFilterQuery(filterData, filterOptions);
          if (filterQuery) {
            params.set('filterQuery', filterQuery);
          }
        }
      }

      const queryString = params.toString();
      return `${URL_PREFIX}/api/${BFF_API_VERSION}/model_catalog/models/export${queryString ? `?${queryString}` : ''}`;
    },
    [queryParams, selectedSourceLabel, searchTerm, filterData, filterOptions],
  );

  const handleExportSelected = React.useCallback(() => {
    setIsOpen(false);
    const models = Array.from(selectedModels.values());
    const csv = generateCatalogModelCSV(models);
    downloadCSV(csv, getExportFilename());
    clearSelection();
  }, [selectedModels, clearSelection]);

  const handleExportFiltered = React.useCallback(() => {
    setIsOpen(false);
    setIsExporting(true);
    const url = buildExportUrl(true);
    triggerBFFDownload(url).finally(() => setIsExporting(false));
  }, [buildExportUrl]);

  const handleExportAll = React.useCallback(() => {
    setIsOpen(false);
    setIsExporting(true);
    const url = buildExportUrl(false);
    triggerBFFDownload(url).finally(() => setIsExporting(false));
  }, [buildExportUrl]);

  const selectedCount = selectedModels.size;

  return (
    <Dropdown
      isOpen={isOpen}
      onSelect={() => setIsOpen(false)}
      onOpenChange={setIsOpen}
      toggle={(toggleRef: React.Ref<MenuToggleElement>) => (
        <MenuToggle
          ref={toggleRef}
          onClick={() => setIsOpen(!isOpen)}
          isExpanded={isOpen}
          isDisabled={isExporting}
          variant="secondary"
          data-testid="export-catalog-button"
          icon={isExporting ? <Spinner size="sm" /> : <FileExportIcon />}
        >
          Export
        </MenuToggle>
      )}
    >
      <DropdownList>
        <DropdownItem
          key="export-selected"
          onClick={handleExportSelected}
          isDisabled={selectedCount === 0}
          data-testid="export-selected-option"
          description="Export only the models you have selected"
        >
          {selectedCount > 0 ? `Export selected (${selectedCount})` : 'Export selected'}
        </DropdownItem>
        <DropdownItem
          key="export-filtered"
          onClick={handleExportFiltered}
          data-testid="export-filtered-option"
          description="Export models matching current filters"
        >
          Export current view
        </DropdownItem>
        <DropdownItem
          key="export-all"
          onClick={handleExportAll}
          data-testid="export-all-option"
          description="Export all models in the catalog"
        >
          Export all models
        </DropdownItem>
      </DropdownList>
    </Dropdown>
  );
};

const triggerBFFDownload = async (url: string): Promise<void> => {
  const response = await fetch(url);
  if (!response.ok) {
    throw new Error(`Export failed: ${response.statusText}`);
  }
  const blob = await response.blob();
  const blobUrl = URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = blobUrl;

  const disposition = response.headers.get('Content-Disposition');
  const filenameMatch = disposition?.match(/filename="?([^"]+)"?/);
  link.download = filenameMatch?.[1] ?? getExportFilename();

  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  URL.revokeObjectURL(blobUrl);
};

export default ModelCatalogExportButton;
