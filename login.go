package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ajananias/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerAuth(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		User
		Token string `json:"token"`
	}
	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds *int   `json:"expires_in_seconds"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	}
	db_hashedPassword, err := cfg.db.CheckDBHash(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	match, err := auth.CheckPasswordHash(params.Password, db_hashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to check password", err)
		return
	}
	if !match {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	db_user, err := cfg.db.GetUserFromEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	var expiresIn time.Duration
	if params.ExpiresInSeconds == nil || *params.ExpiresInSeconds >= 3600 {
		expiresIn = 3600 * time.Second
	} else {
		expiresIn = time.Duration(*params.ExpiresInSeconds) * time.Second
	}
	tokenString, err := auth.MakeJWT(db_user.ID, cfg.tokenSecret, expiresIn)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate token for authentication", err)
		return
	}

	user := User{
		ID:        db_user.ID,
		CreatedAt: db_user.CreatedAt,
		UpdatedAt: db_user.UpdatedAt,
		Email:     db_user.Email,
	}
	response := Response{
		User:  user,
		Token: tokenString,
	}
	respondWithJSON(w, http.StatusOK, response)
}
