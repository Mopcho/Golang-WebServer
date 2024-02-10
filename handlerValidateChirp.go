package main

import (
	"encoding/json"
	"net/http"
)

func (apiCnfg *apiCnfg) handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type bodyParameters struct {
		Body string `json:"body"`
	}

	params := bodyParameters{}

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&params)

	if err != nil {
		type errorMsg struct {
			Error string `json:"error"`
		}

		errorMsgStruct := errorMsg{
			Error: "Something went wrong",
		}

		respondWithJSON(w, r, errorMsgStruct, 400)
		return
	}

	if len(params.Body) > 140 {
		type errorMsg struct {
			Error string `json:"error"`
		}

		errorMsgStruct := errorMsg{
			Error: "Chirp is too long",
		}

		respondWithJSON(w, r, errorMsgStruct, 400)
		return
	}

	type validMsg struct {
		Valid bool `json:"valid"`
	}

	validMsgParams := validMsg{
		Valid: true,
	}

	respondWithJSON(w, r, validMsgParams, 200)
}
