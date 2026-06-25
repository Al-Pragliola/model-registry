export const AGENT_CATALOG_SETTINGS_PAGE_TITLE = 'Agent catalog settings';
export const AGENT_CATALOG_SETTINGS_DESCRIPTION =
  'Add and manage agent sources that populate the agent catalog for users in your organization.';

export const ADD_AGENT_SOURCE_TITLE = 'Add a source';
export const ADD_AGENT_SOURCE_DESCRIPTION = 'Add a new agent catalog source to your organization.';

export const MANAGE_AGENT_SOURCE_TITLE = 'Manage source';
export const MANAGE_AGENT_SOURCE_DESCRIPTION = 'Manage the selected agent catalog source.';

export const agentCatalogSettingsUrl = (): string => '/agent-catalog-settings';

export const addAgentSourceUrl = (): string => `${agentCatalogSettingsUrl()}/add-source`;

export const manageAgentSourceUrl = (catalogSourceId: string): string =>
  `${agentCatalogSettingsUrl()}/manage-source/${encodeURIComponent(catalogSourceId)}`;
