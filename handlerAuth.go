package main

import (
	"encoding/json"
	"fmt"
	"internal/database"
	"net/http"
)

type UserLoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (apiCnfg *apiCnfg) handleLogin(w http.ResponseWriter, r *http.Request) {
	params := UserLoginData{}

	decoder := json.NewDecoder(r.Body)

	decoder.Decode(&params)

	// Find user by email
	userInDb, err := database.GetUserByEmail(params.Email)

	if err != nil {
		respondWithError(w, r, fmt.Sprintf("Wrong credentials"), 401)
		return
	}

	// Compare it to the email pass
	err = database.ComparePassword(userInDb.Password, params.Password)

	if err != nil {
		respondWithError(w, r, fmt.Sprintf("Wrong credentials"), 401)
		return
	}

	respondWithJSON(w, r, struct{}{}, 200)
}
