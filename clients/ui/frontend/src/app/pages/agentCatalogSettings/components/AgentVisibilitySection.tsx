import * as React from 'react';
import {
  FormFieldGroupExpandable,
  FormFieldGroupHeader,
  TextArea,
  FormHelperText,
  HelperText,
  HelperTextItem,
} from '@patternfly/react-core';
import { UpdateObjectAtPropAndValue, ThemeAwareFormGroupWrapper } from 'mod-arch-shared';
import FormSection from '~/app/pages/modelRegistry/components/pf-overrides/FormSection';
import { ManageAgentSourceFormData } from '~/app/pages/agentCatalogSettings/useManageAgentSourceData';

type AgentVisibilitySectionProps = {
  formData: ManageAgentSourceFormData;
  setData: UpdateObjectAtPropAndValue<ManageAgentSourceFormData>;
  isDefaultExpanded?: boolean;
};

const AgentVisibilitySection: React.FC<AgentVisibilitySectionProps> = ({
  formData,
  setData,
  isDefaultExpanded = false,
}) => {
  const includedAgentsInput = (
    <TextArea
      id="included-agents"
      name="included-agents"
      data-testid="included-agents-input"
      value={formData.includedAgents}
      onChange={(_event, value) => setData('includedAgents', value)}
      rows={3}
      resizeOrientation="vertical"
      placeholder="Example: conversational-*, tool-use-agent"
    />
  );

  const includedAgentsDescriptionTxtNode = (
    <FormHelperText>
      <HelperText>
        <HelperTextItem>
          Enter the names of agents to include from this source. These agents will appear in the
          agent catalog.
        </HelperTextItem>
      </HelperText>
    </FormHelperText>
  );

  const includedAgentsHelperTxtNode = (
    <FormHelperText>
      <HelperText>
        <HelperTextItem>
          Separate agent names using commas. To include all agents with a specific prefix, enter the
          prefix followed by an asterisk. Example: conversational-*
        </HelperTextItem>
      </HelperText>
    </FormHelperText>
  );

  const excludedAgentsInput = (
    <TextArea
      id="excluded-agents"
      name="excluded-agents"
      data-testid="excluded-agents-input"
      value={formData.excludedAgents}
      onChange={(_event, value) => setData('excludedAgents', value)}
      rows={3}
      resizeOrientation="vertical"
      placeholder="Example: conversational-*, tool-use-agent"
    />
  );

  const excludedAgentsDescriptionTxtNode = (
    <FormHelperText>
      <HelperText>
        <HelperTextItem>
          Enter the names of agents to exclude from this source. These agents will not appear in the
          agent catalog.
        </HelperTextItem>
      </HelperText>
    </FormHelperText>
  );

  const excludedAgentsHelperTxtNode = (
    <FormHelperText>
      <HelperText>
        <HelperTextItem>
          Separate agent names using commas. To exclude all agents with a specific prefix, enter the
          prefix followed by an asterisk. Example: conversational-*
        </HelperTextItem>
      </HelperText>
    </FormHelperText>
  );

  return (
    <FormSection>
      <FormFieldGroupExpandable
        toggleAriaLabel="Agent visibility"
        header={
          <FormFieldGroupHeader
            titleText={{ text: 'Agent visibility', id: 'agent-visibility-title' }}
            titleDescription="Optionally filter which agents from this source appear in the agent catalog. If no filters are set, all agents from the source will be visible."
          />
        }
        isExpanded={isDefaultExpanded}
        data-testid="agent-visibility-section"
      >
        <ThemeAwareFormGroupWrapper
          label="Included agents"
          fieldId="included-agents"
          descriptionTextNode={includedAgentsDescriptionTxtNode}
          helperTextNode={includedAgentsHelperTxtNode}
        >
          {includedAgentsInput}
        </ThemeAwareFormGroupWrapper>

        <ThemeAwareFormGroupWrapper
          label="Excluded agents"
          fieldId="excluded-agents"
          descriptionTextNode={excludedAgentsDescriptionTxtNode}
          helperTextNode={excludedAgentsHelperTxtNode}
        >
          {excludedAgentsInput}
        </ThemeAwareFormGroupWrapper>
      </FormFieldGroupExpandable>
    </FormSection>
  );
};

export default AgentVisibilitySection;
