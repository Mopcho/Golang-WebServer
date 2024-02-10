package main

import (
	"encoding/json"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, r *http.Request, data interface{}, statusCode int) {
	jsonData, err := json.Marshal(data)

	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(statusCode)
	w.Write(jsonData)
}
