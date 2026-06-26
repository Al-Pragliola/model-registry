package models

type AgentEnvVar struct {
	Name     string `json:"name"`
	Required bool   `json:"required"`
}

type AgentArtifact struct {
	URI                      string  `json:"uri"`
	CreateTimeSinceEpoch     *string `json:"createTimeSinceEpoch,omitempty"`
	LastUpdateTimeSinceEpoch *string `json:"lastUpdateTimeSinceEpoch,omitempty"`
}

type Agent struct {
	ID                       string          `json:"id"`
	Name                     string          `json:"name"`
	SourceID                 *string         `json:"source_id,omitempty"`
	DisplayName              *string         `json:"displayName,omitempty"`
	Description              *string         `json:"description,omitempty"`
	Framework                *string         `json:"framework,omitempty"`
	AgentType                *string         `json:"agentType,omitempty"`
	Tags                     []string        `json:"tags,omitempty"`
	Models                   []string        `json:"models,omitempty"`
	Logo                     *string         `json:"logo,omitempty"`
	RepositoryUrl            *string         `json:"repositoryUrl,omitempty"`
	PublishedDate            *string         `json:"publishedDate,omitempty"`
	Readme                   *string         `json:"readme,omitempty"`
	Env                      []AgentEnvVar   `json:"env,omitempty"`
	Artifacts                []AgentArtifact `json:"artifacts,omitempty"`
	CreateTimeSinceEpoch     *string         `json:"createTimeSinceEpoch,omitempty"`
	LastUpdateTimeSinceEpoch *string         `json:"lastUpdateTimeSinceEpoch,omitempty"`
}

type AgentList struct {
	NextPageToken string  `json:"nextPageToken"`
	PageSize      int32   `json:"pageSize"`
	Size          int32   `json:"size"`
	Items         []Agent `json:"items"`
}
