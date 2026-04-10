package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/AggroSec/Go-HTTP-Server/internal/auth"
	"github.com/AggroSec/Go-HTTP-Server/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) chirpHandler(w http.ResponseWriter, r *http.Request) {
	type SentChirp struct {
		Body   string `json:"body"`
		UserID string `json:"user_id"`
	}

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Authorization header not found")
	}

	userID, err := auth.ValidateJWT(tokenString, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	decoder := json.NewDecoder(r.Body)
	chirp := SentChirp{}
	err = decoder.Decode(&chirp)
	if err != nil {
		respondWithError(w, http.StatusExpectationFailed, "Something went wrong")
		return
	}

	chirp.UserID = userID.String()

	if len(chirp.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	chirp.Body = cleanBody(chirp.Body)

	savedChirp, err := cfg.db.AddChirp(r.Context(), database.AddChirpParams{
		Body:   chirp.Body,
		UserID: uuid.MustParse(chirp.UserID),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	postedChirp := Chirp{
		ID:        savedChirp.ID,
		CreatedAt: savedChirp.CreatedAt,
		UpdatedAt: savedChirp.UpdatedAt,
		Body:      savedChirp.Body,
		UserID:    savedChirp.UserID,
	}

	respondWithJSON(w, http.StatusCreated, postedChirp)
}

func cleanBody(body string) string {
	split_body := strings.Split(body, " ")
	bad_words := map[string]bool{
		"kerfuffle": true,
		"sharbert":  true,
		"fornax":    true,
	}
	clean_body_slice := []string{}

	for _, word := range split_body {
		lowerWord := strings.ToLower(word)
		if bad_words[lowerWord] {
			clean_body_slice = append(clean_body_slice, "****")
		} else {
			clean_body_slice = append(clean_body_slice, word)
		}
	}
	clean_body := strings.Join(clean_body_slice, " ")
	return clean_body
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	authorID := r.URL.Query().Get("author_id")
	sorting := r.URL.Query().Get("sort")

	if sorting != "" && sorting != "asc" && sorting != "desc" {
		respondWithError(w, http.StatusBadRequest, "Invalid sort parameter")
		return
	}

	if authorID != "" {
		authorUUID, err := uuid.Parse(authorID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author_id")
			return
		}
		chirpsByAuthor, err := cfg.GetChirpsByAuthor(authorUUID, r)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Something went wrong")
			return
		}
		switch sorting {
		case "asc":
			sort.Slice(chirpsByAuthor, func(i, j int) bool {
				return chirpsByAuthor[i].CreatedAt.Before(chirpsByAuthor[j].CreatedAt)
			})
		case "desc":
			sort.Slice(chirpsByAuthor, func(i, j int) bool {
				return chirpsByAuthor[i].CreatedAt.After(chirpsByAuthor[j].CreatedAt)
			})
		}
		respondWithJSON(w, http.StatusOK, chirpsByAuthor)
		return
	}

	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	chirpsRetrieved := []Chirp{}
	for _, chirp := range chirps {
		chirpsRetrieved = append(chirpsRetrieved, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	switch sorting {
	case "asc":
		sort.Slice(chirpsRetrieved, func(i, j int) bool {
			return chirpsRetrieved[i].CreatedAt.Before(chirpsRetrieved[j].CreatedAt)
		})
	case "desc":
		sort.Slice(chirpsRetrieved, func(i, j int) bool {
			return chirpsRetrieved[i].CreatedAt.After(chirpsRetrieved[j].CreatedAt)
		})
	}

	respondWithJSON(w, http.StatusOK, chirpsRetrieved)
}

func (cfg *apiConfig) GetChirpsByAuthor(authorID uuid.UUID, r *http.Request) ([]Chirp, error) {
	chirps, err := cfg.db.GetChirpsByAuthor(r.Context(), authorID)
	if err != nil {
		return nil, err
	}

	chirpsRetrieved := []Chirp{}
	for _, chirp := range chirps {
		chirpsRetrieved = append(chirpsRetrieved, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	return chirpsRetrieved, nil
}

func (cfg *apiConfig) handlerGetChirpByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("chirpID")
	chirp, err := cfg.db.GetChirpByID(r.Context(), uuid.MustParse(id))
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found")
		return
	}

	chirpRetrieved := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	respondWithJSON(w, http.StatusOK, chirpRetrieved)
}

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Authorization header not found")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	chirpID := r.PathValue("chirpID")

	chirp, err := cfg.db.GetChirpByID(r.Context(), uuid.MustParse(chirpID))
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found")
		return
	}

	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "You can only delete your own chirps")
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), database.DeleteChirpParams{
		ID:     uuid.MustParse(chirpID),
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
