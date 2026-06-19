package agentcatalog

import (
	agentArtifactmodels "github.com/kubeflow/hub/catalog/internal/catalog/agentcatalog/models"
	agentmodels "github.com/kubeflow/hub/catalog/internal/catalog/agentcatalog/models"
	sharedmodels "github.com/kubeflow/hub/catalog/internal/db/models"
)

type Services struct {
	AgentRepository           agentmodels.AgentRepository
	AgentArtifactRepository   agentArtifactmodels.AgentArtifactRepository
	CatalogSourceRepository   sharedmodels.CatalogSourceRepository
	PropertyOptionsRepository sharedmodels.PropertyOptionsRepository
}
