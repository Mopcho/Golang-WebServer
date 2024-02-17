package main

import (
	"encoding/json"
	"internal/database"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserLoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (apiCnfg *apiCnfg) handleLogin(w http.ResponseWriter, r *http.Request) {
	params := UserLoginData{}

	decoder := json.NewDecoder(r.Body)

	decoder.Decode(&params)

	userInDb, err := database.GetUserByEmail(params.Email)

	if err != nil {
		respondWithError(w, r, "Wrong credentials", 401)
		return
	}

	err = database.ComparePassword(userInDb.Password, params.Password)

	if err != nil {
		respondWithError(w, r, "Wrong credentials", 401)
		return
	}

	// Login success
	tokenString, err := apiCnfg.createToken(userInDb.ID)

	if err != nil {
		respondWithError(w, r, "Could not sign JWT", 500)
		return
	}

	type Response struct {
		JWT string `json:"jwt"`
	}

	respondWithJSON(w, r, Response{JWT: tokenString}, 200)
}

func (apiCnfg *apiCnfg) createToken(userId string) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		Subject:   userId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(apiCnfg.jwtSecret))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (apiCnfg *apiCnfg) keyFunc(token *jwt.Token) (interface{}, error) {
	return []byte(apiCnfg.jwtSecret), nil
}

func (apiCnfg *apiCnfg) DecodeToken(token string) (jwt.Claims, error) {
	claims := &jwt.RegisteredClaims{}

	decodedToken, err := jwt.ParseWithClaims(token, claims, apiCnfg.keyFunc)

	if err != nil {
		return nil, err
	}

	return decodedToken.Claims, nil
}
