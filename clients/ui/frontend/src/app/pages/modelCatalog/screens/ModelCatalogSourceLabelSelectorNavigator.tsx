import React from 'react';
import { CatalogModel } from '~/app/modelCatalogTypes';
import ModelCatalogSourceLabelSelector from './ModelCatalogSourceLabelSelector';

type ModelCatalogSourceLabelSelectorNavigatorProps = {
  searchTerm?: string;
  onSearch?: (term: string) => void;
  onClearSearch?: () => void;
  onResetAllFilters?: () => void;
  selectedModels?: Map<string, CatalogModel>;
  clearSelection?: () => void;
};

const ModelCatalogSourceLabelSelectorNavigator: React.FC<
  ModelCatalogSourceLabelSelectorNavigatorProps
> = ({ searchTerm, onSearch, onClearSearch, onResetAllFilters, selectedModels, clearSelection }) => (
  <ModelCatalogSourceLabelSelector
    searchTerm={searchTerm}
    onSearch={onSearch}
    onClearSearch={onClearSearch}
    onResetAllFilters={onResetAllFilters}
    selectedModels={selectedModels}
    clearSelection={clearSelection}
  />
);
export default ModelCatalogSourceLabelSelectorNavigator;
