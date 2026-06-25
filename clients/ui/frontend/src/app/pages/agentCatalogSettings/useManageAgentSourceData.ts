import { GenericObjectState } from 'mod-arch-core';
import useGenericObjectState from 'mod-arch-core/dist/utilities/useGenericObjectState';

export type ManageAgentSourceFormData = {
  name: string;
  id: string;
  // YAML field
  yamlContent: string;
  // Filter fields
  includedAgents: string;
  excludedAgents: string;
  // Enable source
  enabled: boolean;
  isDefault: boolean;
};

const manageAgentSourceFormDataDefaults: ManageAgentSourceFormData = {
  name: '',
  id: '',
  yamlContent: '',
  includedAgents: '',
  excludedAgents: '',
  enabled: false,
  isDefault: false,
};

/**
 * Custom hook to manage form state for adding/editing an agent catalog source.
 * Agent sources are always YAML-only (no HuggingFace option).
 * Uses the standard useGenericObjectState pattern from mod-arch-core.
 * @param existingData - Optional existing data to pre-populate the form (for edit mode)
 * @returns Generic object state with [formData, setData] pattern
 */
export const useManageAgentSourceData = (
  existingData?: Partial<ManageAgentSourceFormData>,
): GenericObjectState<ManageAgentSourceFormData> =>
  useGenericObjectState<ManageAgentSourceFormData>({
    ...manageAgentSourceFormDataDefaults,
    ...existingData,
  });
