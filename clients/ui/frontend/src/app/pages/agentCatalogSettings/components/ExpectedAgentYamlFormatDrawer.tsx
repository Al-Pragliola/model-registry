import * as React from 'react';
import {
  CodeBlock,
  CodeBlockCode,
  DrawerActions,
  DrawerCloseButton,
  DrawerHead,
  DrawerPanelBody,
} from '@patternfly/react-core';
import sampleAgentCatalogYamlContent from '~/app/pages/agentCatalogSettings/sample-agent-catalog.yaml';

type ExpectedAgentYamlFormatDrawerPanelProps = {
  onClose: () => void;
};

export const ExpectedAgentYamlFormatDrawerPanel: React.FC<
  ExpectedAgentYamlFormatDrawerPanelProps
> = ({ onClose }) => (
  <>
    <DrawerHead>
      <span data-testid="expected-agent-format-drawer-title">View expected file format</span>
      <DrawerActions>
        <DrawerCloseButton
          onClose={onClose}
          aria-label="Close drawer"
          data-testid="expected-agent-format-drawer-close"
        />
      </DrawerActions>
    </DrawerHead>
    <DrawerPanelBody hasNoPadding>
      <CodeBlock>
        <CodeBlockCode>{sampleAgentCatalogYamlContent}</CodeBlockCode>
      </CodeBlock>
    </DrawerPanelBody>
  </>
);
