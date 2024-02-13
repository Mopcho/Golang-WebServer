package main

import (
	"github.com/go-chi/chi/v5"
)

func setupApiRouter(apiRouter *chi.Mux, apiCnfg *apiCnfg) {
	apiRouter.Get("/healthz", handleHealthz)
	apiRouter.Get("/metrics", apiCnfg.handleMetricts)
	apiRouter.Get("/metrics/reset", apiCnfg.handleMetrictsReset)
	apiRouter.Post("/chirps", apiCnfg.handleCreateChirp)
	apiRouter.Get("/chirps/{chirpId}", apiCnfg.handlerGetChirp)
	apiRouter.Get("/chirps", apiCnfg.handlerGetChirps)

	apiRouter.Post("/users", apiCnfg.handlerCreateUser)
	apiRouter.Get("/users/{userId}", apiCnfg.handlerGetUser)
	apiRouter.Get("/users", apiCnfg.handleGetUsers)
	apiRouter.Post("/login", apiCnfg.handleLogin)
}
