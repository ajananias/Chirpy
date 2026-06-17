package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ajananias/Chirpy/internal/auth"
	"github.com/ajananias/Chirpy/internal/database"

	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
		//UserID uuid.UUID `json:"user_id"`
	}

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "401 Unauthorized", err)
		return
	}
	validUserID, err := auth.ValidateJWT(bearerToken, cfg.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Failed to validate token", err)
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	maxChirpLength := 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}
	clean_body := cleanBody(params.Body)
	respBody := validResponse{
		CleanBody: clean_body,
	}

	db_chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   respBody.CleanBody,
		UserID: validUserID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create chirp", err)
		return
	}

	chirp := Chirp{
		ID:        db_chirp.ID,
		CreatedAt: db_chirp.CreatedAt,
		UpdatedAt: db_chirp.UpdatedAt,
		Body:      db_chirp.Body,
		UserID:    db_chirp.UserID,
	}

	respondWithJSON(w, http.StatusCreated, chirp)

}
