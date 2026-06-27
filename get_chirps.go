package main

import (
	"net/http"
	"sort"

	"github.com/ajananias/Chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	var db_chirps []database.Chirp
	queryParams := r.URL.Query()
	stringAuthorID := queryParams.Get("author_id")
	authorID, err := uuid.Parse(stringAuthorID)
	if err == nil {
		dbUserChirps, err := cfg.db.GetChirpsByUserID(r.Context(), authorID)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "Couldn't find user from provided id", err)
			return
		}
		db_chirps = dbUserChirps
	} else {
		db_chirps, err = cfg.db.GetAllChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Unable to retrieve chirps", err)
			return
		}
	}

	if queryParams.Has("sort") {
		chirpSort := queryParams.Get("sort")
		switch chirpSort {
		case "asc":
			sort.Slice(db_chirps, func(i, j int) bool { return db_chirps[i].CreatedAt.Before(db_chirps[j].CreatedAt) })
		case "desc":
			sort.Slice(db_chirps, func(i, j int) bool { return db_chirps[i].CreatedAt.After(db_chirps[j].CreatedAt) })
		}
	} else {
		sort.Slice(db_chirps, func(i, j int) bool { return db_chirps[i].CreatedAt.Before(db_chirps[j].CreatedAt) })
	}

	all_chirps := make([]Chirp, len(db_chirps))
	for i, db_chirp := range db_chirps {
		all_chirps[i] = Chirp{
			ID:        db_chirp.ID,
			CreatedAt: db_chirp.CreatedAt,
			UpdatedAt: db_chirp.UpdatedAt,
			Body:      db_chirp.Body,
			UserID:    db_chirp.UserID,
		}
	}
	respondWithJSON(w, http.StatusOK, all_chirps)
}

func (cfg *apiConfig) handlerGetChirpByID(w http.ResponseWriter, r *http.Request) {
	raw_chirp_id := r.PathValue("chirpID")
	chirp_id, err := uuid.Parse(raw_chirp_id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}
	db_chirp, err := cfg.db.GetChirpByChirpID(r.Context(), chirp_id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}
	chirp := Chirp{
		ID:        db_chirp.ID,
		CreatedAt: db_chirp.CreatedAt,
		UpdatedAt: db_chirp.UpdatedAt,
		Body:      db_chirp.Body,
		UserID:    db_chirp.UserID,
	}
	respondWithJSON(w, http.StatusOK, chirp)
}
