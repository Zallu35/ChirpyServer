package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Zallu35/ChirpyServer/internal/auth"
	"github.com/google/uuid"
)

type refreshResponse struct {
	Token string `json:"token"`
}

func (a *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Refresh failed: %v", err)
		errorResponse(w, http.StatusUnauthorized, "No Authorization header found")
		return
	}
	user, err := a.database.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		log.Printf("Refresh token error: %v", err)
		errorResponse(w, http.StatusInternalServerError, "Error fetching user data")
		return
	}
	if user.UserID == uuid.Nil {
		log.Printf("Refresh token does not exist")
		errorResponse(w, http.StatusUnauthorized, "Invalid token")
		return
	}
	if user.ExpiresAt.Before(time.Now()) {
		log.Printf("Refresh token expired")
		errorResponse(w, http.StatusUnauthorized, "Token expired")
		return
	}
	if user.RevokedAt.Valid != false && user.RevokedAt.Time.Before(time.Now()) {
		log.Printf("Refresh token revoked")
		errorResponse(w, http.StatusUnauthorized, "Token revoked")
		return
	}
	accToken, err := auth.MakeJWT(user.UserID, a.secret)
	if err != nil {
		log.Printf("Refresh - Error making new JWT: %v", err)
		errorResponse(w, http.StatusInternalServerError, "Error refreshing token")
		return
	}
	rr := refreshResponse{
		Token: accToken,
	}
	respondWithJSON(w, http.StatusOK, rr)
}

func (a *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Revoke failed: %v", err)
		errorResponse(w, http.StatusUnauthorized, "No Authorization header found")
		return
	}
	er := a.database.RevokeToken(r.Context(), refreshToken)
	if er != nil {
		log.Printf("Revoke failed: %v", er)
		errorResponse(w, http.StatusInternalServerError, "Error revoking token")
	}

	w.WriteHeader(http.StatusNoContent)
}
