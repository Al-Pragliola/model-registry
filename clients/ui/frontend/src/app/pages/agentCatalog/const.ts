import type { AgentFilterCategoryKey } from '~/app/pages/agentCatalog/types/agentCatalogFilterOptions';

export const AGENT_CATALOG_TITLE = 'Agent Catalog';
export const AGENT_CATALOG_DESCRIPTION =
  'Browse and deploy AI agents built with popular frameworks.';

export const AGENT_CATALOG_GALLERY = {
  CARDS_PER_ROW: 4,
  PAGE_SIZE: 10,
  SECTION_TITLE: 'Agents',
} as const;

type GridSpan = 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 | 11 | 12;

const GRID_COLUMNS = 12;
const GRID_SPAN_VALUES: GridSpan[] = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12];

function toGridSpan(cols: number): GridSpan {
  const index = Math.min(Math.max(0, cols - 1), GRID_SPAN_VALUES.length - 1);
  return GRID_SPAN_VALUES[index];
}

export const AGENT_CATALOG_GRID_SPAN: {
  sm: GridSpan;
  md: GridSpan;
  lg: GridSpan;
  xl2: GridSpan;
} = {
  sm: toGridSpan(GRID_COLUMNS),
  md: toGridSpan(GRID_COLUMNS / 2),
  lg: toGridSpan(GRID_COLUMNS / AGENT_CATALOG_GALLERY.CARDS_PER_ROW),
  xl2: toGridSpan(GRID_COLUMNS / AGENT_CATALOG_GALLERY.CARDS_PER_ROW),
};

export const AGENT_FILTER_CATEGORY_NAMES: Record<AgentFilterCategoryKey, string> = {
  framework: 'Framework',
  agentType: 'Agent type',
  labels: 'Labels',
};

export const AGENT_FILTER_KEYS: AgentFilterCategoryKey[] = ['framework', 'agentType', 'labels'];

export const BACKEND_TO_FRONTEND_AGENT_FILTER_KEY: Record<string, AgentFilterCategoryKey> = {
  tags: 'labels',
};

export const OTHER_AGENTS_DISPLAY_NAME = 'Other agents';
