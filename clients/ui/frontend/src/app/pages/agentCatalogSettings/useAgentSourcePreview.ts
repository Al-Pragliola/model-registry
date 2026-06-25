import * as React from 'react';
import { isEqual } from 'lodash-es';
import {
  CatalogSourceConfig,
  CatalogSourcePreviewRequest,
  CatalogSourcePreviewModel,
  CatalogSourcePreviewSummary,
} from '~/app/modelCatalogTypes';
import { AgentCatalogSettingsAPIState } from '~/app/hooks/agentCatalogSettings/useAgentCatalogSettingsAPIState';
import { ManageAgentSourceFormData } from './useManageAgentSourceData';

export enum PreviewMode {
  PREVIEW = 'preview',
  VALIDATE = 'validate',
}

export enum PreviewTab {
  INCLUDED = 'included',
  EXCLUDED = 'excluded',
}

const DEFAULT_PREVIEW_PAGE_SIZE = 20;

const getTargetTab = (
  isFreshPreview: boolean,
  switchToTab: PreviewTab | undefined,
  activeTab: PreviewTab,
): PreviewTab => {
  if (isFreshPreview) {
    return PreviewTab.INCLUDED;
  }
  return switchToTab ?? activeTab;
};

export type PreviewTabState = {
  items: CatalogSourcePreviewModel[];
  nextPageToken?: string;
  hasMore: boolean;
};

const initialTabState: PreviewTabState = {
  items: [],
  nextPageToken: undefined,
  hasMore: false,
};

export type PreviewState = {
  mode?: PreviewMode;
  isLoadingInitial: boolean;
  isLoadingMore: boolean;
  summary?: CatalogSourcePreviewSummary;
  tabStates: Record<PreviewTab, PreviewTabState>;
  error?: Error;
  resultDismissed: boolean;
  lastPreviewedData?: CatalogSourcePreviewRequest;
  activeTab: PreviewTab;
};

export interface UseAgentSourcePreviewOptions {
  formData: ManageAgentSourceFormData;
  existingSourceConfig?: CatalogSourceConfig;
  apiState: AgentCatalogSettingsAPIState;
  isEditMode: boolean;
}

export interface UseAgentSourcePreviewResult {
  previewState: PreviewState;
  handlePreview: (mode?: PreviewMode) => Promise<void>;
  handleTabChange: (tab: PreviewTab) => void;
  handleLoadMore: () => void;
  hasFormChanged: boolean;
  canPreview: boolean;
}

const parseCommaSeparatedList = (value: string): string[] =>
  value
    .split(',')
    .map((item) => item.trim())
    .filter((item) => item.length > 0);

export const useAgentSourcePreview = ({
  formData,
  existingSourceConfig,
  apiState,
  isEditMode,
}: UseAgentSourcePreviewOptions): UseAgentSourcePreviewResult => {
  const [previewState, setPreviewState] = React.useState<PreviewState>({
    isLoadingInitial: false,
    isLoadingMore: false,
    tabStates: {
      [PreviewTab.INCLUDED]: initialTabState,
      [PreviewTab.EXCLUDED]: initialTabState,
    },
    resultDismissed: false,
    activeTab: PreviewTab.INCLUDED,
  });

  const canPreview = formData.yamlContent.trim().length > 0 || !!formData.isDefault;

  const buildPreviewRequest = React.useCallback((): CatalogSourcePreviewRequest => {
    const request: CatalogSourcePreviewRequest = {
      type: 'yaml',
      includedModels: parseCommaSeparatedList(formData.includedAgents),
      excludedModels: parseCommaSeparatedList(formData.excludedAgents),
    };

    request.properties = {
      yaml: formData.yamlContent,
      yamlCatalogPath:
        existingSourceConfig?.type === 'yaml' ? existingSourceConfig.yamlCatalogPath : undefined,
    };

    return request;
  }, [formData, existingSourceConfig]);

  const hasFormChanged = React.useMemo(() => {
    if (!previewState.lastPreviewedData) {
      return false;
    }
    const currentRequest = buildPreviewRequest();
    return !isEqual(currentRequest, previewState.lastPreviewedData);
  }, [buildPreviewRequest, previewState.lastPreviewedData]);

  const handlePreview = React.useCallback(
    async (
      mode: PreviewMode = PreviewMode.PREVIEW,
      options?: {
        loadMore?: boolean;
        switchToTab?: PreviewTab;
      },
    ) => {
      const { loadMore = false, switchToTab } = options ?? {};
      const isFreshPreview = !loadMore && !switchToTab;
      const targetTab = getTargetTab(isFreshPreview, switchToTab, previewState.activeTab);

      if (!apiState.apiAvailable) {
        setPreviewState((prev) => ({
          ...prev,
          mode,
          isLoadingInitial: false,
          error: new Error('API is not available'),
          resultDismissed: false,
        }));
        return;
      }

      if (isFreshPreview) {
        setPreviewState({
          mode,
          isLoadingInitial: true,
          isLoadingMore: false,
          tabStates: {
            [PreviewTab.INCLUDED]: initialTabState,
            [PreviewTab.EXCLUDED]: initialTabState,
          },
          activeTab: PreviewTab.INCLUDED,
          error: undefined,
          resultDismissed: false,
          summary: undefined,
          lastPreviewedData: undefined,
        });
      } else if (loadMore) {
        setPreviewState((prev) => ({ ...prev, isLoadingMore: true }));
      } else if (switchToTab) {
        setPreviewState((prev) => ({ ...prev, activeTab: switchToTab, isLoadingInitial: true }));
      }

      let requestData: CatalogSourcePreviewRequest;
      if (isFreshPreview) {
        requestData = buildPreviewRequest();
      } else if (previewState.lastPreviewedData) {
        requestData = previewState.lastPreviewedData;
      } else {
        return handlePreview(mode);
      }

      const nextPageToken = loadMore ? previewState.tabStates[targetTab].nextPageToken : undefined;

      try {
        const result = await apiState.api.previewCatalogSource({}, requestData, {
          filterStatus: targetTab,
          pageSize: DEFAULT_PREVIEW_PAGE_SIZE,
          nextPageToken,
        });

        setPreviewState((prev) => {
          const currentTabState = prev.tabStates[targetTab];
          const newItems = loadMore ? [...currentTabState.items, ...result.items] : result.items;

          return {
            ...prev,
            mode,
            isLoadingInitial: false,
            isLoadingMore: false,
            summary: result.summary,
            lastPreviewedData: isFreshPreview ? requestData : prev.lastPreviewedData,
            tabStates: {
              ...prev.tabStates,
              [targetTab]: {
                items: newItems,
                nextPageToken: result.nextPageToken,
                hasMore: !!result.nextPageToken && result.items.length > 0,
              },
            },
            error: undefined,
            resultDismissed: false,
          };
        });
      } catch (error) {
        const err = error instanceof Error ? error : new Error('Failed to preview source');

        setPreviewState((prev) => ({
          ...prev,
          mode,
          isLoadingInitial: false,
          isLoadingMore: false,
          error: err,
          resultDismissed: false,
        }));
      }
    },
    [
      apiState,
      buildPreviewRequest,
      previewState.activeTab,
      previewState.lastPreviewedData,
      previewState.tabStates,
    ],
  );

  const handleTabChange = React.useCallback(
    (newTab: PreviewTab) => {
      if (newTab === previewState.activeTab) {
        return;
      }

      const tabState = previewState.tabStates[newTab];
      if (tabState.items.length === 0) {
        handlePreview(PreviewMode.PREVIEW, { switchToTab: newTab });
      } else {
        setPreviewState((prev) => ({ ...prev, activeTab: newTab }));
      }
    },
    [handlePreview, previewState.activeTab, previewState.tabStates],
  );

  const handleLoadMore = React.useCallback(() => {
    handlePreview(PreviewMode.PREVIEW, { loadMore: true });
  }, [handlePreview]);

  // Auto-trigger preview on mount in edit mode
  React.useEffect(() => {
    const hasNoResults = previewState.tabStates[PreviewTab.INCLUDED].items.length === 0;
    if (isEditMode && canPreview && hasNoResults) {
      handlePreview();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return {
    previewState,
    handlePreview,
    handleTabChange,
    handleLoadMore,
    hasFormChanged,
    canPreview,
  };
};
