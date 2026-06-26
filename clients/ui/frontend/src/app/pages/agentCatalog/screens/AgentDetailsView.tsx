import * as React from 'react';
import {
  Card,
  CardBody,
  CardHeader,
  ClipboardCopy,
  Content,
  DescriptionList,
  DescriptionListDescription,
  DescriptionListGroup,
  DescriptionListTerm,
  Icon,
  Label,
  LabelGroup,
  PageSection,
  Sidebar,
  SidebarContent,
  SidebarPanel,
  Stack,
  StackItem,
  Title,
} from '@patternfly/react-core';
import { GithubIcon, OutlinedClockIcon } from '@patternfly/react-icons';
import type { Agent } from '~/app/agentCatalogTypes';
import ExternalLink from '~/app/shared/components/ExternalLink';
import MarkdownComponent from '~/app/shared/markdown/MarkdownComponent';
import ModelTimestamp from '~/app/pages/modelRegistry/screens/components/ModelTimestamp';

type AgentDetailsViewProps = {
  agent: Agent;
};

const VISIBLE_LABELS = 5;

const AgentDetailsView: React.FC<AgentDetailsViewProps> = ({ agent }) => (
  <PageSection hasBodyWrapper={false} isFilled padding={{ default: 'noPadding' }}>
    <Sidebar hasGutter isPanelRight>
      <SidebarContent style={{ minWidth: 0, overflow: 'hidden' }}>
        <Stack hasGutter>
          <StackItem>
            <Card>
              <CardHeader>
                <Title headingLevel="h2" size="lg">
                  Description
                </Title>
              </CardHeader>
              <CardBody>
                <Content className="pf-v6-u-text-break-word">
                  <p data-testid="agent-description">
                    {agent.description || 'No description'}
                  </p>
                </Content>
              </CardBody>
            </Card>
          </StackItem>
          {agent.env && agent.env.length > 0 && (
            <StackItem>
              <Card>
                <CardHeader>
                  <Title headingLevel="h2" size="lg">
                    Environment Variables
                  </Title>
                </CardHeader>
                <CardBody>
                  <DescriptionList isHorizontal>
                    {agent.env.map((envVar) => (
                      <DescriptionListGroup key={envVar.name}>
                        <DescriptionListTerm>
                          <code>{envVar.name}</code>
                        </DescriptionListTerm>
                        <DescriptionListDescription>
                          {envVar.required ? (
                            <Label color="red" isCompact>
                              Required
                            </Label>
                          ) : (
                            <Label color="grey" isCompact>
                              Optional
                            </Label>
                          )}
                        </DescriptionListDescription>
                      </DescriptionListGroup>
                    ))}
                  </DescriptionList>
                </CardBody>
              </Card>
            </StackItem>
          )}
          <StackItem>
            <Card>
              <CardHeader>
                <Title headingLevel="h2" size="lg">
                  <Icon isInline style={{ marginRight: '4px' }}>
                    <GithubIcon />
                  </Icon>
                  README
                </Title>
              </CardHeader>
              <CardBody>
                {!agent.readme && (
                  <Content component="p" data-testid="agent-no-readme">
                    No README available
                  </Content>
                )}
                {agent.readme && (
                  <MarkdownComponent
                    data={agent.readme}
                    dataTestId="agent-readme-markdown"
                    maxHeading={3}
                  />
                )}
              </CardBody>
            </Card>
          </StackItem>
        </Stack>
      </SidebarContent>
      <SidebarPanel width={{ default: 'width_33' }}>
        <Card>
          <CardHeader>
            <Title headingLevel="h2" size="lg">
              Agent details
            </Title>
          </CardHeader>
          <CardBody>
            <DescriptionList>
              {agent.tags && agent.tags.length > 0 && (
                <DescriptionListGroup>
                  <DescriptionListTerm>Labels</DescriptionListTerm>
                  <DescriptionListDescription>
                    <LabelGroup numLabels={VISIBLE_LABELS} isCompact>
                      {agent.tags.map((tag) => (
                        <Label key={tag} variant="outline" data-testid="agent-detail-label">
                          {tag}
                        </Label>
                      ))}
                    </LabelGroup>
                  </DescriptionListDescription>
                </DescriptionListGroup>
              )}
              {agent.framework && (
                <DescriptionListGroup>
                  <DescriptionListTerm>Framework</DescriptionListTerm>
                  <DescriptionListDescription data-testid="agent-framework">
                    {agent.framework}
                  </DescriptionListDescription>
                </DescriptionListGroup>
              )}
              {agent.agentType && (
                <DescriptionListGroup>
                  <DescriptionListTerm>Agent type</DescriptionListTerm>
                  <DescriptionListDescription data-testid="agent-type">
                    {agent.agentType}
                  </DescriptionListDescription>
                </DescriptionListGroup>
              )}
              {agent.models && agent.models.length > 0 && (
                <DescriptionListGroup>
                  <DescriptionListTerm>Compatible models</DescriptionListTerm>
                  <DescriptionListDescription>
                    <Stack>
                      {agent.models.map((model) => (
                        <StackItem key={model} data-testid="agent-compatible-model">
                          {model}
                        </StackItem>
                      ))}
                    </Stack>
                  </DescriptionListDescription>
                </DescriptionListGroup>
              )}
              {agent.artifacts && agent.artifacts.length > 0 && (
                <DescriptionListGroup>
                  <DescriptionListTerm>Artifacts</DescriptionListTerm>
                  <DescriptionListDescription>
                    <Stack hasGutter>
                      {agent.artifacts.map((artifact) => (
                        <StackItem key={artifact.uri}>
                          <ClipboardCopy
                            hoverTip="Copy"
                            clickTip="Copied"
                            isReadOnly
                            data-testid="agent-artifact-copy"
                          >
                            {artifact.uri}
                          </ClipboardCopy>
                        </StackItem>
                      ))}
                    </Stack>
                  </DescriptionListDescription>
                </DescriptionListGroup>
              )}
              {agent.repositoryUrl && (
                <DescriptionListGroup>
                  <DescriptionListTerm>Repository</DescriptionListTerm>
                  <DescriptionListDescription>
                    <ExternalLink
                      text={agent.repositoryUrl}
                      to={agent.repositoryUrl}
                      testId="agent-repo-link"
                    />
                  </DescriptionListDescription>
                </DescriptionListGroup>
              )}
              {agent.publishedDate && (
                <DescriptionListGroup>
                  <DescriptionListTerm>Published</DescriptionListTerm>
                  <DescriptionListDescription>
                    <Icon isInline style={{ marginRight: '4px' }}>
                      <OutlinedClockIcon />
                    </Icon>
                    <ModelTimestamp timeSinceEpoch={agent.publishedDate} />
                  </DescriptionListDescription>
                </DescriptionListGroup>
              )}
            </DescriptionList>
          </CardBody>
        </Card>
      </SidebarPanel>
    </Sidebar>
  </PageSection>
);

export default AgentDetailsView;
