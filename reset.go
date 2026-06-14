package main

import (
	"net/http"
	"os"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0\n"))

	PLATFORM := os.Getenv("PLATFORM")
	if PLATFORM != "dev" {
		http.Error(w, "403 Forbidden", http.StatusForbidden)
		return
	}
	err := cfg.db.ResetUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't reset users", err)
		return
	}
	w.Write([]byte("All users deleted successfully\n"))

}
