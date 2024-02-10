package main

import (
	"net/http"
)

func respondWithError(w http.ResponseWriter, r *http.Request, msg string, statusCode int) {
	type errorMsg struct {
		Error string `json:"error"`
	}

	errorMsgStruct := errorMsg{
		Error: msg,
	}

	respondWithJSON(w, r, errorMsgStruct, statusCode)
}
