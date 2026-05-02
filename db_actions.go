package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (a *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	type userEmail struct {
		Email string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	requestData := userEmail{}
	err := decoder.Decode(&requestData)
	if err != nil {
		log.Printf("createUser - Error decoding request: %s", err)
		errorResponse(w, http.StatusInternalServerError, "Error decoding request")
		return
	}

	usr, err := a.database.CreateUser(r.Context(), requestData.Email)
	if err != nil {
		log.Printf("createUser - Error creating user: %s", err)
		errorResponse(w, http.StatusInternalServerError, "Error creating user")
		return
	}

	jsonUsr := User{
		usr.ID,
		usr.CreatedAt,
		usr.UpdatedAt,
		usr.Email,
	}
	respondWithJSON(w, http.StatusCreated, jsonUsr)
}
