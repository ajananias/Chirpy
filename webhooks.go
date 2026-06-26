package main

import (
	"encoding/json"
	"net/http"

	"github.com/ajananias/Chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerWebhooks(w http.ResponseWriter, r *http.Request) {
	type DataParameters struct {
		UserID uuid.UUID `json:"user_id"`
	}
	type Parameters struct {
		Event string         `json:"event"`
		Data  DataParameters `json:"data"`
	}
	type Response struct {
		Body string `json:"body"`
	}
	apiKey, err := auth.GetAPIKey(r.Header)
	if apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Wrong API key", err)
	}
	decoder := json.NewDecoder(r.Body)
	params := Parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding webhook parameters", err)
		return
	}
	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
	}
	userID := params.Data.UserID
	err = cfg.db.UpgradeMembership(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't find & upgrade user", err)
		return
	}
	response := Response{
		Body: "",
	}
	respondWithJSON(w, http.StatusNoContent, response)
}
