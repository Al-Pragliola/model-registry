package models

import (
	dbmodels "github.com/kubeflow/hub/internal/platform/db/entity"
	"github.com/kubeflow/hub/internal/platform/db/filter"
)

// AgentArtifactListOptions holds the options for listing AgentArtifact entities.
type AgentArtifactListOptions struct {
	dbmodels.Pagination
	SourceIDs   *[]string
	FilterQuery *string
}

// GetRestEntityType implements the FilterApplier interface.
func (o *AgentArtifactListOptions) GetRestEntityType() filter.RestEntityType {
	return filter.RestEntityType("agent_artifact")
}

// GetFilterQuery returns the filter query string for advanced filtering.
func (o *AgentArtifactListOptions) GetFilterQuery() string {
	if o.FilterQuery == nil {
		return ""
	}
	return *o.FilterQuery
}

// AgentArtifactAttributes holds the attributes for a AgentArtifact record.
type AgentArtifactAttributes struct {
	Name                     *string
	ExternalID               *string
	CreateTimeSinceEpoch     *int64
	LastUpdateTimeSinceEpoch *int64
}

// AgentArtifact represents a AgentArtifact stored in the database.
type AgentArtifact interface {
	dbmodels.Entity[AgentArtifactAttributes]
}

// AgentArtifactImpl is the concrete implementation of AgentArtifact.
type AgentArtifactImpl = dbmodels.BaseEntity[AgentArtifactAttributes]

// AgentArtifactRepository defines the interface for AgentArtifact persistence.
type AgentArtifactRepository interface {
	GetByID(id int32) (AgentArtifact, error)
	GetByName(name string) (AgentArtifact, error)
	List(listOptions *AgentArtifactListOptions) (*dbmodels.ListWrapper[AgentArtifact], error)
	Save(entity AgentArtifact, parentResourceID *int32) (AgentArtifact, error)
	DeleteBySource(sourceID string) error
	DeleteByID(id int32) error
	GetDistinctSourceIDs() ([]string, error)
	GetTypeID() int32
}
