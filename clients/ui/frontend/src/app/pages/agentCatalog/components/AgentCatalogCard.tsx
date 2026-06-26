import * as React from 'react';
import {
  Button,
  Card,
  CardBody,
  CardHeader,
  CardTitle,
  Flex,
  FlexItem,
  Label,
  LabelGroup,
  Truncate,
} from '@patternfly/react-core';
import { TruncatedText } from 'mod-arch-shared';
import { AutomationIcon } from '@patternfly/react-icons';
import { Link, type LinkProps } from 'react-router-dom';
import type { Agent } from '~/app/agentCatalogTypes';
import { agentDetailsUrl } from '~/app/routes/agentCatalog/agentCatalog';

type AgentCatalogCardProps = {
  agent: Agent;
};

const AgentCatalogCard: React.FC<AgentCatalogCardProps> = React.memo(({ agent }) => {
  const agentId = agent.id;
  const displayTitle = agent.displayName || agent.name;

  return (
    <Card isFullHeight data-testid={`agent-catalog-card-${agentId}`}>
      <CardHeader>
        <Flex
          alignItems={{ default: 'alignItemsFlexStart' }}
          justifyContent={{ default: 'justifyContentSpaceBetween' }}
          gap={{ default: 'gapXs' }}
          className="pf-v6-u-mb-md"
        >
          <FlexItem>
            {agent.logo ? (
              <img
                src={agent.logo}
                alt=""
                style={{ height: '32px', width: '32px' }}
                data-testid={`agent-catalog-card-logo-${agentId}`}
              />
            ) : (
              <span
                className="pf-v6-u-display-inline-block pf-v6-u-font-size-2xl pf-v6-u-color-200"
                aria-hidden
              >
                <AutomationIcon />
              </span>
            )}
          </FlexItem>
          {agent.framework && (
            <FlexItem>
              <Label data-testid={`agent-catalog-card-framework-${agentId}`}>
                {agent.framework}
              </Label>
            </FlexItem>
          )}
        </Flex>
        <CardTitle>
          <Button
            data-testid={`agent-catalog-card-detail-link-${agentId}`}
            variant="link"
            isInline
            component={(props: LinkProps) => <Link {...props} to={agentDetailsUrl(agentId)} />}
            style={{
              fontSize: 'var(--pf-t--global--font--size--body--default)',
              fontWeight: 'var(--pf-t--global--font--weight--body--bold)',
            }}
          >
            <Truncate
              content={displayTitle}
              position="middle"
              tooltipPosition="top"
              data-testid={`agent-catalog-card-name-${agentId}`}
            />
          </Button>
        </CardTitle>
      </CardHeader>
      <CardBody>
        <TruncatedText
          content={agent.description ?? ''}
          maxLines={4}
          data-testid={`agent-catalog-card-description-${agentId}`}
        />
        {agent.tags && agent.tags.length > 0 && (
          <LabelGroup numLabels={3} isCompact className="pf-v6-u-mt-lg">
            {agent.tags.map((tag) => (
              <Label key={tag} variant="outline">
                {tag}
              </Label>
            ))}
          </LabelGroup>
        )}
      </CardBody>
    </Card>
  );
});
AgentCatalogCard.displayName = 'AgentCatalogCard';

export default AgentCatalogCard;
