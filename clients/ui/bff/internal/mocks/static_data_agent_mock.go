package mocks

import "github.com/kubeflow/hub/ui/bff/internal/models"

func GetAgentMocks() []models.Agent {
	sourceID := "agent-source-1"
	framework1 := "langgraph"
	framework2 := "crewai"
	framework3 := "a2a"
	agentType := "starter_kit"
	desc1 := "A conversational AI agent built with LangGraph for multi-turn dialogue and tool use."
	desc2 := "A collaborative crew of agents for automated code review and analysis."
	desc3 := "An A2A-compatible agent for document processing and summarization."
	displayName1 := "Conversational Assistant"
	displayName2 := "Code Review Crew"
	displayName3 := "Document Processor"
	repoURL1 := "https://github.com/example/conversational-assistant"
	repoURL2 := "https://github.com/example/code-review-crew"
	repoURL3 := "https://github.com/example/doc-processor"
	readme := "# Agent README\n\nThis is a sample agent readme with **markdown** support."

	return []models.Agent{
		{
			ID:          "1",
			Name:        "conversational-assistant",
			SourceID:    &sourceID,
			DisplayName: &displayName1,
			Description: &desc1,
			Framework:   &framework1,
			AgentType:   &agentType,
			Tags:        []string{"conversational", "tool-use", "multi-turn"},
			Models:      []string{"granite-3b", "granite-8b"},
			RepositoryUrl: &repoURL1,
			Readme:      &readme,
			Env: []models.AgentEnvVar{
				{Name: "OPENAI_API_KEY", Required: true},
				{Name: "LOG_LEVEL", Required: false},
			},
			Artifacts: []models.AgentArtifact{
				{URI: "oci://registry.example.com/agents/conversational-assistant:v1.0"},
			},
		},
		{
			ID:          "2",
			Name:        "code-review-crew",
			SourceID:    &sourceID,
			DisplayName: &displayName2,
			Description: &desc2,
			Framework:   &framework2,
			AgentType:   &agentType,
			Tags:        []string{"code-review", "automation", "devtools"},
			Models:      []string{"granite-8b-code-instruct"},
			RepositoryUrl: &repoURL2,
			Readme:      &readme,
			Env: []models.AgentEnvVar{
				{Name: "GITHUB_TOKEN", Required: true},
			},
			Artifacts: []models.AgentArtifact{
				{URI: "oci://registry.example.com/agents/code-review-crew:v2.1"},
			},
		},
		{
			ID:          "3",
			Name:        "doc-processor",
			SourceID:    &sourceID,
			DisplayName: &displayName3,
			Description: &desc3,
			Framework:   &framework3,
			AgentType:   &agentType,
			Tags:        []string{"documents", "summarization", "a2a"},
			Models:      []string{"granite-3b"},
			RepositoryUrl: &repoURL3,
			Readme:      &readme,
			Env:         []models.AgentEnvVar{},
			Artifacts: []models.AgentArtifact{
				{URI: "oci://registry.example.com/agents/doc-processor:v1.2"},
			},
		},
	}
}

func GetAgentFilterOptionsListMock() models.FilterOptionsList {
	frameworkOption := models.FilterOption{
		Type: "string",
		Values: []interface{}{
			"langgraph",
			"crewai",
			"autogen",
			"a2a",
		},
	}
	agentTypeOption := models.FilterOption{
		Type: "string",
		Values: []interface{}{
			"starter_kit",
		},
	}
	tagsOption := models.FilterOption{
		Type: "string",
		Values: []interface{}{
			"conversational",
			"tool-use",
			"code-review",
			"automation",
			"documents",
			"a2a",
		},
	}

	filters := map[string]models.FilterOption{
		"framework": frameworkOption,
		"agentType": agentTypeOption,
		"tags":      tagsOption,
	}

	return models.FilterOptionsList{
		Filters: &filters,
	}
}
