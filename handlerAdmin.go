package main

import (
	"fmt"
	"net/http"
)

func (apiCnfg *apiCnfg) handleAdminPage(w http.ResponseWriter, r *http.Request) {
	pageContent := generateAdminPage(apiCnfg.fileserverHits)

	w.WriteHeader(200)
	w.Header().Add("COntent-Type", "text/html")
	w.Write([]byte(pageContent))
}

func generateAdminPage(visitedTimes int) string {
	return fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", visitedTimes)
}
