package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"slices"

	"github.com/google/uuid"
)

type Chirp struct {
	ID   string `json:"id"`
	Body string `json:"body"`
}

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

	chirp := Chirp{
		Body: validChirp,
		ID:   uuid.NewString(),
	}

	err = saveChirpToDisk(chirp)

	if err != nil {
		respondWithError(w, r, "Failed saving file to disk", 500)
		return
	}
	respondWithJSON(w, r, chirp, 201)
}

func saveChirpToDisk(chirp Chirp) error {
	f, err := os.OpenFile("./chirps.txt", os.O_CREATE, 0660)
	defer f.Close()

	readBytes, err := os.ReadFile("./chirps.txt")

	if err != nil {
		return errors.New("Something went wrong")
	}

	readChirps := make([]Chirp, 0)

	err = json.Unmarshal(readBytes, &readChirps)

	if err != nil {
		return errors.New("Something went wrong")
	}

	newChirps := slices.Insert(readChirps, 0, chirp)

	newChirpsBytes, err := json.Marshal(newChirps)

	if err != nil {
		return errors.New("Something went wrong")
	}

	_, err = f.Write(newChirpsBytes)
	return err
}
