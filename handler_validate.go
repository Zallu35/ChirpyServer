package main

import (
	"encoding/json"
	"log"
	"net/http"
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
		Valid bool `json:"valid"`
	}

	if len(msg.Message) > 140 {
		errorResponse(w, http.StatusBadRequest, "Chipt is too long")
		return
	}

	rsp := validationReply{
		Valid: true,
	}
	respondWithJSON(w, http.StatusOK, rsp)
}
