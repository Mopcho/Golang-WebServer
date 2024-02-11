package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"internal/database"

	"github.com/go-chi/chi/v5"
)

func (apiCnfg *apiCnfg) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type bodyParameters struct {
		Email string `json:"email"`
	}

	params := bodyParameters{}

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, r, "Something Went wrong", 400)
		return
	}

	user := database.UserCreateData{
		Email: params.Email,
	}

	err = database.CreateUser(user)

	if err != nil {
		respondWithError(w, r, "Failed saving file to disk", 500)
		return
	}
	respondWithJSON(w, r, user, 201)
}

func (apiCnfg *apiCnfg) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := database.GetUsers()

	if err != nil {
		respondWithError(w, r, fmt.Sprintf("Error getting users: %v", err), 500)
		return
	}

	respondWithJSON(w, r, users, 200)
}

func (apiCnfg *apiCnfg) handlerGetUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")

	user, err := database.GetUser(userId)

	if err != nil {
		respondWithError(w, r, fmt.Sprintf("Error getting user: %v", err), 500)
		return
	}

	respondWithJSON(w, r, user, 200)
}
