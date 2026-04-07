package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func chirpValidationHandler(w http.ResponseWriter, r *http.Request) {
	type Chirp struct {
		Body string `json:"body"`
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

	chirp.Body = cleanBody(chirp.Body)

	respondWithJSON(w, 200, map[string]string{"cleaned_body": chirp.Body})
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
