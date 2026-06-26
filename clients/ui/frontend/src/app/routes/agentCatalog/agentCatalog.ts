export const agentCatalogUrl = (): string => '/agent-catalog';

export const agentDetailsUrl = (agentId: string | number): string =>
  `${agentCatalogUrl()}/${encodeURIComponent(String(agentId))}`;
