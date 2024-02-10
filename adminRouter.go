package main

import "github.com/go-chi/chi/v5"

func setupAdminRouter(adminRouter *chi.Mux, apiCnfg *apiCnfg) {
	adminRouter.Get("/metrics", apiCnfg.handleAdminPage)
}
