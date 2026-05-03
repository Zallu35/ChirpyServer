package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/Zallu35/ChirpyServer/internal/auth"
	"github.com/Zallu35/ChirpyServer/internal/database"
	"github.com/google/uuid"
)

func (a *apiConfig) handlerPostChirp(w http.ResponseWriter, r *http.Request) {
	type chirpCheck struct {
		Message string `json:"body"`
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting authorization token: %v", err)
		errorResponse(w, http.StatusUnauthorized, "Error getting authorization token")
		return
	}
	userID, err := auth.ValidateJWT(token, a.secret)
	if err != nil {
		log.Printf("Error validating token: %v", err)
		errorResponse(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	decoder := json.NewDecoder(r.Body)
	msg := chirpCheck{}
	er := decoder.Decode(&msg)
	if er != nil {
		log.Printf("Error decoding chirp to validate: %s", er)
		errorResponse(w, http.StatusInternalServerError, "Error decoding chirp to validate")
		return
	}

	if len(msg.Message) > 140 {
		errorResponse(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	postParams := database.PostChirpParams{Body: cleanChirp(msg.Message), UserID: userID}

	pst, err := a.database.PostChirp(r.Context(), postParams)
	if err != nil {
		log.Printf("postChirp - Error posting chirp: %s", err)
		errorResponse(w, http.StatusInternalServerError, "Error posting chirp")
		return
	}

	rsp := Post{
		ID:        pst.ID,
		CreatedAt: pst.CreatedAt,
		UpdatedAt: pst.UpdatedAt,
		Body:      pst.Body,
		UserID:    pst.UserID,
	}
	respondWithJSON(w, http.StatusCreated, rsp)
}

func cleanChirp(original string) string {
	profaneList := []string{"kerfuffle", "sharbert", "fornax"}
	splitString := strings.Split(original, " ")
	var cleanedSlice = []string{}
	for _, word := range splitString {
		if exists := slices.Contains(profaneList, strings.ToLower(word)); exists == true {
			cleanedSlice = append(cleanedSlice, "****")
			continue
		}
		cleanedSlice = append(cleanedSlice, word)
	}
	return strings.Join(cleanedSlice, " ")
}

func (a *apiConfig) handlerRetrieveChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := a.database.RetrieveChirps(r.Context())
	if err != nil {
		log.Printf("RetrieveChirps - Error retrieving chirps: %s", err)
		errorResponse(w, http.StatusInternalServerError, "Error retrieving chirps")
		return
	}
	chirpArray := []Post{}
	for i := range chirps {
		jsonChirp := Post{
			chirps[i].ID,
			chirps[i].CreatedAt,
			chirps[i].UpdatedAt,
			chirps[i].Body,
			chirps[i].UserID,
		}
		chirpArray = append(chirpArray, jsonChirp)
	}
	respondWithJSON(w, http.StatusOK, chirpArray)
}

func (a *apiConfig) handlerRetrieveSingleChirp(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")
	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		log.Printf("RetrieveSingleChirp - Error Parsing ChirpID: %s", err)
		errorResponse(w, http.StatusInternalServerError, "Error Parsing ChirpID")
		return
	}

	chirp, er := a.database.GetSingleChirp(r.Context(), chirpUUID)
	if er != nil {
		log.Printf("RetrieveSingleChirp - Chirp not found: %s", er)
		errorResponse(w, http.StatusNotFound, "Chirp not found")
		return
	}

	jsonChirp := Post{
		chirp.ID,
		chirp.CreatedAt,
		chirp.UpdatedAt,
		chirp.Body,
		chirp.UserID,
	}
	respondWithJSON(w, http.StatusOK, jsonChirp)
}
