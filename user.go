package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ajananias/Chirpy/internal/auth"
	"github.com/ajananias/Chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

func (cfg *apiConfig) handlerUsers(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	}
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't retrieve bearer access token", err)
		return
	}
	userID, err := auth.ValidateJWT(accessToken, cfg.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate user's access token", err)
		return
	}
	newEmail, newPassword := params.Email, params.Password
	hashedNewPassword, err := auth.HashPassword(newPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash new password", err)
		return
	}
	err = cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		Email:          newEmail,
		HashedPassword: hashedNewPassword,
		ID:             userID,
	})
	dbUpdatedUser, err := cfg.db.GetUserFromEmail(r.Context(), newEmail)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve updated user", err)
		return
	}
	user := User{
		ID:          dbUpdatedUser.ID,
		CreatedAt:   dbUpdatedUser.CreatedAt,
		UpdatedAt:   dbUpdatedUser.UpdatedAt,
		Email:       dbUpdatedUser.Email,
		IsChirpyRed: dbUpdatedUser.IsChirpyRed,
	}
	respondWithJSON(w, http.StatusOK, user)
}
