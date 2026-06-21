package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ajananias/Chirpy/internal/auth"
	"github.com/ajananias/Chirpy/internal/database"
)

func (cfg *apiConfig) handlerAuth(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}
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

	accessExpirationTime := time.Hour
	refreshExpirationTime := time.Hour * 24 * 60
	accessToken, err := auth.MakeJWT(db_user.ID, cfg.tokenSecret, accessExpirationTime)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate token for authentication", err)
		return
	}
	refreshToken := auth.MakeRefreshToken()
	_, err = cfg.db.AddRefreshToken(r.Context(), database.AddRefreshTokenParams{
		Token:     refreshToken,
		UserID:    db_user.ID,
		ExpiresAt: time.Now().Add(refreshExpirationTime),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't save refresh token", err)
		return
	}

	user := User{
		ID:        db_user.ID,
		CreatedAt: db_user.CreatedAt,
		UpdatedAt: db_user.UpdatedAt,
		Email:     db_user.Email,
	}
	response := Response{
		User:         user,
		Token:        accessToken,
		RefreshToken: refreshToken,
	}
	respondWithJSON(w, http.StatusOK, response)
}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Token string `json:"token"`
	}
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "401 Unauthorized", err)
		return
	}

	dbUserID, err := cfg.db.GetUserFromRefreshToken(r.Context(), bearerToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "401 Unauthorized", err)
		return
	}
	newAccessToken, err := auth.MakeJWT(dbUserID, cfg.tokenSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create a new access token", err)
		return
	}

	response := Response{
		Token: newAccessToken,
	}
	respondWithJSON(w, http.StatusOK, response)
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "401 Unauthorized", err)
		return
	}
	err = cfg.db.RevokeToken(r.Context(), bearerToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to revoke token", err)
	}
	w.WriteHeader(http.StatusNoContent)
}
