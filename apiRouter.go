package main

import (
	"github.com/go-chi/chi/v5"
)

func setupApiRouter(apiRouter *chi.Mux, apiCnfg *apiCnfg) {
	apiRouter.Get("/healthz", handleHealthz)
	apiRouter.Get("/metrics", apiCnfg.handleMetricts)
	apiRouter.Get("/metrics/reset", apiCnfg.handleMetrictsReset)
	apiRouter.Post("/chirps", apiCnfg.handleCreateChirp)
}
