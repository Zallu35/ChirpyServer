package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
)

func (a *apiConfig) lengthValidationFunc(w http.ResponseWriter, r *http.Request) {
	type chirpCheck struct {
		Message string `json:"body"`
	}
	decoder := json.NewDecoder(r.Body)
	msg := chirpCheck{}
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("Error decoding chirp to validate: %s", err)
		errorResponse(w, http.StatusInternalServerError, "Error decoding chirp to validate")
		return
	}

	type validationReply struct {
		Valid       bool   `json:"valid"`
		CleanedBody string `json:"cleaned_body"`
	}

	if len(msg.Message) > 140 {
		errorResponse(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	rsp := validationReply{
		Valid:       true,
		CleanedBody: cleanChirp(msg.Message),
	}
	respondWithJSON(w, http.StatusOK, rsp)
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
