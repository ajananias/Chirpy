package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
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
	respondWithJSON(w, http.StatusOK, respBody)
}

func cleanBody(body string) string {
	profane_words := []string{"kerfuffle", "sharbert", "fornax"}
	replacement := "****"
	wordList := strings.Split(body, " ")
	for i, word := range wordList {
		for _, profane_word := range profane_words {
			if strings.ToLower(word) == profane_word {
				wordList[i] = replacement
			}
		}
	}
	clean_body := strings.Join(wordList, " ")

	return clean_body
}
