package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/AggroSec/Go-HTTP-Server/internal/auth"
	"github.com/AggroSec/Go-HTTP-Server/internal/database"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type CreateUserRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	var req CreateUserRequest
	err := decoder.Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "something went wrong in decoding")
		return
	}

	if req.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Email is required")
		return
	}

	if req.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Password is required")
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	createdUser, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          req.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusExpectationFailed, "Something went wrong")
		return
	}

	user := User{
		ID:        createdUser.ID,
		CreatedAt: createdUser.CreatedAt,
		UpdatedAt: createdUser.UpdatedAt,
		Email:     createdUser.Email,
	}

	respondWithJSON(w, http.StatusCreated, user)
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	var req LoginRequest
	err := decoder.Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "something went wrong in decoding")
		return
	}

	dbuser, err := cfg.db.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	match, err := auth.CheckPasswordHash(req.Password, dbuser.HashedPassword)
	if err != nil || !match {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	token, err := auth.MakeJWT(dbuser.ID, cfg.jwtSecret, time.Duration(3600)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	refreshToken := auth.MakeRefreshToken()
	_, err = cfg.db.AddRefreshToken(r.Context(), database.AddRefreshTokenParams{
		Token:  refreshToken,
		UserID: dbuser.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	user := User{
		ID:           dbuser.ID,
		CreatedAt:    dbuser.CreatedAt,
		UpdatedAt:    dbuser.UpdatedAt,
		Email:        dbuser.Email,
		Token:        token,
		RefreshToken: refreshToken,
	}

	respondWithJSON(w, http.StatusOK, user)
}

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	type RefreshResponse struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Authorization header format")
		return
	}

	dbToken, err := cfg.db.VerifyRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	if dbToken.ExpiresAt.Before(time.Now()) || dbToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Refresh token is expired or revoked")
		return
	}

	newToken, err := auth.MakeJWT(dbToken.UserID, cfg.jwtSecret, time.Duration(3600)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	respondWithJSON(w, http.StatusOK, RefreshResponse{Token: newToken})
}

func (cfg *apiConfig) handlerRevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Authorization header format")
		return
	}

	err = cfg.db.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (cfg *apiConfig) handlerUpdateUserLogin(w http.ResponseWriter, r *http.Request) {
	type UpdateUserLoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var req UpdateUserLoginRequest
	err = decoder.Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "something went wrong in decoding")
		return
	}

	if req.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Email is required")
		return
	}

	if req.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Password is required")
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	updatedUser, err := cfg.db.UpdateUserLogin(r.Context(), database.UpdateUserLoginParams{
		ID:             userID,
		Email:          req.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusExpectationFailed, "Something went wrong")
		return
	}

	user := User{
		ID:        updatedUser.ID,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
		Email:     updatedUser.Email,
	}

	respondWithJSON(w, http.StatusOK, user)
}
