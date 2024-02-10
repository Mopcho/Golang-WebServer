package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (apiCnfg *apiCnfg) handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type bodyParameters struct {
		Body string `json:"body"`
	}

	params := bodyParameters{}

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, r, "Something went wrong", 400)
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, r, "Chirp is too long", 400)
		return
	}

	censoredSentance := censorWord(params.Body)
	if params.Body != censoredSentance {
		type validMsg struct {
			CleanedBody string `json:"cleaned_body"`
		}

		validMsgParams := validMsg{
			CleanedBody: censoredSentance,
		}

		respondWithJSON(w, r, validMsgParams, 200)
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

func censorWord(stringToCensor string) string {
	stringToCensor = strings.ReplaceAll(stringToCensor, "kerfuffle", "*********")
	stringToCensor = strings.ReplaceAll(stringToCensor, "sharbert", "********")
	stringToCensor = strings.ReplaceAll(stringToCensor, "fornax", "******")

	return stringToCensor
}
