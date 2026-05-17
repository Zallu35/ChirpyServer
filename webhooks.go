package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Zallu35/ChirpyServer/internal/auth"
	"github.com/google/uuid"
)

func (a *apiConfig) updateChirpyRed(w http.ResponseWriter, r *http.Request) {
	requestKey, err := auth.GetAPIKey(r.Header)
	if requestKey != a.polka_key {
		log.Printf("Invalid API Key")
		errorResponse(w, http.StatusUnauthorized, "Invalid Key")
	}

	type polkaRequest struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	polkDat := polkaRequest{}
	err = decoder.Decode(&polkDat)
	if err != nil {
		log.Printf("updateChirpyRed: %v", err)
		errorResponse(w, http.StatusInternalServerError, "Error decdoing request body")
		return
	}

	if polkDat.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	newRed, err := uuid.Parse(polkDat.Data.UserID)
	if err != nil {
		log.Printf("updateChirpyRed: %v", err)
		errorResponse(w, http.StatusInternalServerError, "Error parsing user ID")
		return
	}
	err = a.database.UpgradeChirpyRed(r.Context(), newRed)
	if err != nil {
		log.Printf("updateChirpyRed: %v", err)
		errorResponse(w, http.StatusNotFound, "User ID unknown or invalid")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
