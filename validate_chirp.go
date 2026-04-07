package main

import (
	"encoding/json"
	"net/http"
)

func chirpValidationHandler(w http.ResponseWriter, r *http.Request) {
	type Chirp struct {
		Body string `json:"body"`
	}

	type ValidResponse struct {
		Valid bool `json:"valid"`
	}

	decoder := json.NewDecoder(r.Body)
	chirp := Chirp{}
	err := decoder.Decode(&chirp)
	if err != nil {
		respondWithError(w, http.StatusExpectationFailed, "Something went wrong")
		return
	}

	if len(chirp.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	respondWithJSON(w, 200, ValidResponse{Valid: true})
}
