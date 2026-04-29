package api

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/kubeflow/hub/ui/bff/internal/constants"
	"github.com/kubeflow/hub/ui/bff/internal/integrations/httpclient"
)

func (app *App) ExportCatalogModelsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	client, ok := r.Context().Value(constants.ModelCatalogHttpClientKey).(httpclient.HTTPClientInterface)
	if !ok {
		app.serverErrorResponse(w, r, errors.New("catalog REST client not found"))
		return
	}

	dryRun := r.URL.Query().Get("dryRun") == "true"

	responseData, err := app.repositories.ModelCatalogClient.ExportCatalogModels(client, r.URL.Query())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if dryRun {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(responseData)
		return
	}

	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="model-catalog-export.csv"`)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(responseData)
}
