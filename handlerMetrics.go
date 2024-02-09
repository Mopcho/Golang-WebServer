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
