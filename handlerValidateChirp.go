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

	// Censor it
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
