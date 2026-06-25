package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/kubeflow/hub/ui/bff/internal/constants"
	"github.com/kubeflow/hub/ui/bff/internal/integrations/httpclient"
	"github.com/kubeflow/hub/ui/bff/internal/models"
	"github.com/kubeflow/hub/ui/bff/internal/repositories"
)

type AgentCatalogSettingsSourceConfigEnvelope Envelope[*models.CatalogSourceConfig, None]
type AgentCatalogSettingsSourceConfigListEnvelope Envelope[*models.CatalogSourceConfigList, None]
type AgentCatalogSourcePayloadEnvelope Envelope[*models.CatalogSourceConfigPayload, None]

func (app *App) GetAllAgentCatalogSourceConfigsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()

	namespace, ok := ctx.Value(constants.NamespaceHeaderParameterKey).(string)
	if !ok || namespace == "" {
		app.badRequestResponse(w, r, fmt.Errorf("missing namespace in context"))
		return
	}

	client, err := app.kubernetesClientFactory.GetClient(ctx)
	if err != nil {
		app.serverErrorResponse(w, r, errors.New("catalog client not found"))
		return
	}

	configs, err := app.repositories.AgentCatalogSettingsRepository.GetAllAgentCatalogSourceConfigs(ctx, client, namespace)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	envelope := AgentCatalogSettingsSourceConfigListEnvelope{
		Data: configs,
	}

	if err = app.WriteJSON(w, http.StatusOK, envelope, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *App) GetAgentCatalogSourceConfigHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()

	namespace, ok := ctx.Value(constants.NamespaceHeaderParameterKey).(string)
	if !ok || namespace == "" {
		app.badRequestResponse(w, r, fmt.Errorf("missing namespace in context"))
		return
	}

	catalogSourceId := ps.ByName(CatalogSourceId)

	client, err := app.kubernetesClientFactory.GetClient(ctx)
	if err != nil {
		app.serverErrorResponse(w, r, errors.New("catalog client not found"))
		return
	}

	config, err := app.repositories.AgentCatalogSettingsRepository.GetAgentCatalogSourceConfig(ctx, client, namespace, catalogSourceId)
	if err != nil {
		if errors.Is(err, repositories.ErrCatalogSourceNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	envelope := AgentCatalogSettingsSourceConfigEnvelope{
		Data: config,
	}

	if err = app.WriteJSON(w, http.StatusOK, envelope, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *App) CreateAgentCatalogSourceConfigHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()

	namespace, ok := ctx.Value(constants.NamespaceHeaderParameterKey).(string)
	if !ok || namespace == "" {
		app.badRequestResponse(w, r, fmt.Errorf("missing namespace in context"))
		return
	}

	client, err := app.kubernetesClientFactory.GetClient(ctx)
	if err != nil {
		app.serverErrorResponse(w, r, errors.New("catalog client not found"))
		return
	}

	var envelope AgentCatalogSourcePayloadEnvelope
	if err := json.NewDecoder(r.Body).Decode(&envelope); err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("error decoding JSON: %v", err.Error()))
		return
	}

	newSource, err := app.repositories.AgentCatalogSettingsRepository.CreateAgentCatalogSourceConfig(ctx, client, namespace, *envelope.Data)
	if err != nil {
		if errors.Is(err, repositories.ErrCatalogSourceAlreadyExist) ||
			errors.Is(err, repositories.ErrCatalogSourceIdRequired) ||
			errors.Is(err, repositories.ErrUnsupportedCatalogType) ||
			errors.Is(err, repositories.ErrValidationFailed) {
			app.badRequestResponse(w, r, err)
		} else if errors.Is(err, repositories.ErrCatalogSourceConflict) {
			app.conflictResponse(w, r, err.Error())
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	resultEnvelope := AgentCatalogSettingsSourceConfigEnvelope{
		Data: newSource,
	}

	w.Header().Set("Location", r.URL.JoinPath(resultEnvelope.Data.Id).String())
	if writeErr := app.WriteJSON(w, http.StatusCreated, resultEnvelope, nil); writeErr != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("error writing JSON"))
	}
}

func (app *App) UpdateAgentCatalogSourceConfigHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()

	namespace, ok := ctx.Value(constants.NamespaceHeaderParameterKey).(string)
	if !ok || namespace == "" {
		app.badRequestResponse(w, r, fmt.Errorf("missing namespace in context"))
		return
	}

	client, err := app.kubernetesClientFactory.GetClient(ctx)
	if err != nil {
		app.serverErrorResponse(w, r, errors.New("catalog client not found"))
		return
	}

	var envelope AgentCatalogSourcePayloadEnvelope
	if err := json.NewDecoder(r.Body).Decode(&envelope); err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("error decoding JSON: %v", err.Error()))
		return
	}

	catalogSourceId := ps.ByName(CatalogSourceId)
	if catalogSourceId == "" {
		catalogSourceId = envelope.Data.Id
	}

	updatedSource, err := app.repositories.AgentCatalogSettingsRepository.UpdateAgentCatalogSourceConfig(ctx, client, namespace, catalogSourceId, *envelope.Data)
	if err != nil {
		if errors.Is(err, repositories.ErrCatalogSourceNotFound) {
			app.notFoundResponse(w, r)
		} else if errors.Is(err, repositories.ErrCannotChangeDefaultSource) ||
			errors.Is(err, repositories.ErrCannotChangeType) {
			app.forbiddenResponse(w, r, err.Error())
		} else if errors.Is(err, repositories.ErrCatalogSourceConflict) {
			app.conflictResponse(w, r, err.Error())
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	resultEnvelope := AgentCatalogSettingsSourceConfigEnvelope{
		Data: updatedSource,
	}

	if err = app.WriteJSON(w, http.StatusOK, resultEnvelope, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *App) DeleteAgentCatalogSourceConfigHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()

	namespace, ok := ctx.Value(constants.NamespaceHeaderParameterKey).(string)
	if !ok || namespace == "" {
		app.badRequestResponse(w, r, fmt.Errorf("missing namespace in context"))
		return
	}

	client, err := app.kubernetesClientFactory.GetClient(ctx)
	if err != nil {
		app.serverErrorResponse(w, r, errors.New("catalog client not found"))
		return
	}

	catalogSourceId := ps.ByName(CatalogSourceId)

	deletedSource, err := app.repositories.AgentCatalogSettingsRepository.DeleteAgentCatalogSourceConfig(ctx, client, namespace, catalogSourceId)
	if err != nil {
		if errors.Is(err, repositories.ErrCannotDeleteDefaultSource) {
			app.forbiddenResponse(w, r, err.Error())
		} else if errors.Is(err, repositories.ErrCatalogSourceNotFound) {
			app.notFoundResponse(w, r)
		} else if errors.Is(err, repositories.ErrCatalogSourceConflict) {
			app.conflictResponse(w, r, err.Error())
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	resultEnvelope := AgentCatalogSettingsSourceConfigEnvelope{
		Data: deletedSource,
	}

	if err = app.WriteJSON(w, http.StatusOK, resultEnvelope, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *App) CreateAgentCatalogSourcePreviewHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	client, ok := r.Context().Value(constants.ModelCatalogHttpClientKey).(httpclient.HTTPClientInterface)
	if !ok {
		app.serverErrorResponse(w, r, errors.New("catalog REST client not found"))
		return
	}

	var requestBody struct {
		Data models.CatalogSourcePreviewRequest `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("error decoding JSON: %v", err.Error()))
		return
	}

	previewResult, err := app.repositories.ModelCatalogClient.CreateAgentSourcePreview(client, requestBody.Data, r.URL.Query())
	if err != nil {
		var httpErr *httpclient.HTTPError
		if errors.As(err, &httpErr) {
			app.errorResponse(w, r, httpErr)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	envelope := CatalogSourcePreviewEnvelope{
		Data: previewResult,
	}

	if err = app.WriteJSON(w, http.StatusOK, envelope, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
