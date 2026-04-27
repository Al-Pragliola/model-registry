import React from 'react';
import { CatalogModel } from '~/app/modelCatalogTypes';

const getModelKey = (model: CatalogModel): string => `${model.source_id ?? ''}/${model.name}`;

type UseModelCatalogSelection = {
  selectedModels: Map<string, CatalogModel>;
  selectedCount: number;
  isModelSelected: (model: CatalogModel) => boolean;
  toggleModelSelection: (model: CatalogModel) => void;
  selectAll: (models: CatalogModel[]) => void;
  deselectAll: () => void;
  clearSelection: () => void;
};

export const useModelCatalogSelection = (): UseModelCatalogSelection => {
  const [selectedModels, setSelectedModels] = React.useState<Map<string, CatalogModel>>(
    () => new Map(),
  );

  const isModelSelected = React.useCallback(
    (model: CatalogModel): boolean => selectedModels.has(getModelKey(model)),
    [selectedModels],
  );

  const toggleModelSelection = React.useCallback((model: CatalogModel) => {
    setSelectedModels((prev) => {
      const next = new Map(prev);
      const key = getModelKey(model);
      if (next.has(key)) {
        next.delete(key);
      } else {
        next.set(key, model);
      }
      return next;
    });
  }, []);

  const selectAll = React.useCallback((models: CatalogModel[]) => {
    setSelectedModels((prev) => {
      const next = new Map(prev);
      for (const model of models) {
        next.set(getModelKey(model), model);
      }
      return next;
    });
  }, []);

  const deselectAll = React.useCallback(() => {
    setSelectedModels(new Map());
  }, []);

  return {
    selectedModels,
    selectedCount: selectedModels.size,
    isModelSelected,
    toggleModelSelection,
    selectAll,
    deselectAll,
    clearSelection: deselectAll,
  };
};
