package main

import (
	"fmt"
	"net/http"
)

func (apiCnfg *apiCnfg) handleMetricts(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "plain/text; charset=utf-8")
	w.WriteHeader(200)
	str := fmt.Sprintf("Hits: %v", apiCnfg.fileserverHits)
	w.Write([]byte(str))
}

func (apiCnfg *apiCnfg) handleMetrictsReset(w http.ResponseWriter, r *http.Request) {
	apiCnfg.fileserverHits = 0

	w.WriteHeader(200)
}

func (apiCnfg *apiCnfg) middlewareMetricsInc(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiCnfg.fileserverHits++
		next.ServeHTTP(w, r)
	}
}
