package agent

import (
	"context"

	"github.com/go-chi/chi/v5"

	mapset "github.com/deckarep/golang-set/v2"

	"github.com/kubeflow/hub/catalog/internal/catalog/agentcatalog"
	agentArtifactmodels "github.com/kubeflow/hub/catalog/internal/catalog/agentcatalog/models"
	agentmodels "github.com/kubeflow/hub/catalog/internal/catalog/agentcatalog/models"
	agentservice "github.com/kubeflow/hub/catalog/internal/catalog/agentcatalog/service"
	"github.com/kubeflow/hub/catalog/internal/catalog/basecatalog"
	"github.com/kubeflow/hub/catalog/internal/db/models"
	"github.com/kubeflow/hub/catalog/internal/plugin"
	"github.com/kubeflow/hub/catalog/internal/server/openapi"
	"github.com/kubeflow/hub/internal/platform/datastore"
)

type Plugin struct {
	*plugin.PluginBase
	loader   *agentcatalog.AgentLoader
	services agentcatalog.Services
}

func (p *Plugin) Name() string                   { return "agent" }
func (p *Plugin) Version() string                { return "v1alpha1" }
func (p *Plugin) Description() string            { return "Agent catalog" }
func (p *Plugin) BasePath() string               { return "/api/agent_catalog/v1alpha1" }
func (p *Plugin) Healthy() bool                  { return true }
func (p *Plugin) Migrations() []plugin.Migration { return nil }

func (p *Plugin) DatastoreEntries() []plugin.DatastoreEntry {
	return []plugin.DatastoreEntry{
		{
			TypeName: "kf.Agent",
			Category: "context",
			Spec: datastore.NewSpecType(agentservice.NewAgentRepository).
				AddString("source_id").
				AddString("description").
				AddString("displayName").
				AddString("framework").
				AddString("agentType").
				AddStruct("tags").
				AddStruct("models").
				AddString("logo").
				AddString("repositoryUrl").
				AddString("publishedDate").
				AddString("readme").
				AddStruct("env").
				AddStruct("artifacts"),
		},
		{
			TypeName: "kf.AgentArtifact",
			Category: "artifact",
			Spec: datastore.NewSpecType(agentservice.NewAgentArtifactRepository).
				AddString("source_id"),
		},
	}
}

func (p *Plugin) Init(_ context.Context, cfg plugin.Config) error {
	p.services = agentcatalog.Services{
		AgentRepository:           plugin.GetRepo[agentmodels.AgentRepository](cfg.RepoSet),
		AgentArtifactRepository:   plugin.GetRepo[agentArtifactmodels.AgentArtifactRepository](cfg.RepoSet),
		CatalogSourceRepository:   plugin.GetRepo[models.CatalogSourceRepository](cfg.RepoSet),
		PropertyOptionsRepository: plugin.GetRepo[models.PropertyOptionsRepository](cfg.RepoSet),
	}

	base := basecatalog.NewBaseLoader(cfg.ConfigPaths)
	p.loader = agentcatalog.NewAgentLoader(p.services, base)

	p.PluginBase = plugin.NewPluginBase(plugin.PluginBaseConfig{
		Name:        "agent",
		State:       base,
		Loader:      p.loader,
		FileWatcher: basecatalog.GetMonitor(),
		SourceIDs: func() mapset.Set[string] {
			ids := mapset.NewSet[string]()
			for id := range p.loader.Sources.AllSources() {
				ids.Add(id)
			}
			return ids
		},
	})

	return nil
}

func (p *Plugin) AgentSources() *agentcatalog.AgentSourceCollection {
	return p.loader.Sources
}

func (p *Plugin) RegisterRoutes(router chi.Router) error {
	provider := agentcatalog.NewDBAgentCatalog(p.services, p.loader.Sources)
	svc := openapi.NewAgentCatalogServiceAPIService(provider, p.loader.Sources)
	ctrl := openapi.NewAgentCatalogServiceAPIController(svc)

	for _, route := range ctrl.OrderedRoutes() {
		router.Method(route.Method, route.Pattern, route.HandlerFunc)
	}

	return nil
}
