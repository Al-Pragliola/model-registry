import { PaginationParams } from '~/app/shared/types/catalogTypes';

export type AgentEnvVar = {
  name: string;
  required: boolean;
};

export type AgentArtifact = {
  uri: string;
  createTimeSinceEpoch?: string;
  lastUpdateTimeSinceEpoch?: string;
};

export type Agent = {
  id: string;
  name: string;
  source_id?: string;
  displayName?: string;
  description?: string;
  framework?: string;
  agentType?: string;
  tags?: string[];
  models?: string[];
  logo?: string;
  repositoryUrl?: string;
  publishedDate?: string;
  readme?: string;
  env?: AgentEnvVar[];
  artifacts?: AgentArtifact[];
  createTimeSinceEpoch?: string;
  lastUpdateTimeSinceEpoch?: string;
};

export type AgentList = PaginationParams & { items?: Agent[] };

export type AgentListParams = {
  sourceLabel?: string;
  pageSize?: number | string;
  nextPageToken?: string;
  filterQuery?: string;
  orderBy?: string;
  sortOrder?: string;
  name?: string;
  q?: string;
};
