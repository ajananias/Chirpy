package main

import (
	"encoding/json"
	"net/http"

	"github.com/ajananias/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerAuth(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
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
	user := User{
		ID:        db_user.ID,
		CreatedAt: db_user.CreatedAt,
		UpdatedAt: db_user.UpdatedAt,
		Email:     db_user.Email,
	}
	respondWithJSON(w, http.StatusOK, user)
}
