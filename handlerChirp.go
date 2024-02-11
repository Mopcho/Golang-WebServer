package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"internal/database"
)

func (apiCnfg *apiCnfg) handleCreateChirp(w http.ResponseWriter, r *http.Request) {
	type bodyParameters struct {
		Body string `json:"body"`
	}

	params := bodyParameters{}

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, r, "Something Went wrong", 400)
		return
	}

	validChirp, err := validateChirp(params.Body)

	if err != nil {
		respondWithError(w, r, "Invalid Chirp", 400)
	}

	chirp := database.CreateChirpData{
		Body: validChirp,
	}

	err = database.SaveChirpToDisk(chirp)

	if err != nil {
		respondWithError(w, r, "Failed saving file to disk", 500)
		return
	}
	respondWithJSON(w, r, chirp, 201)
}

func (apiCnfg *apiCnfg) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := database.GetChirpsFromDisk()

	if err != nil {
		respondWithError(w, r, fmt.Sprintf("Error getting chirps: %v", err), 500)
		return
	}

	respondWithJSON(w, r, chirps, 200)
}
