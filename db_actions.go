package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Zallu35/ChirpyServer/internal/auth"
	"github.com/Zallu35/ChirpyServer/internal/database"
)

func (a *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	type userData struct {
		Pw    string `json:"password"`
		Email string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	requestData := userData{}
	err := decoder.Decode(&requestData)
	if err != nil {
		log.Printf("createUser - Error decoding request: %s", err)
		errorResponse(w, http.StatusInternalServerError, "Error decoding request")
		return
	}

	hash, er := auth.HashPassword(requestData.Pw)
	if er != nil {
		log.Printf("createUser - Error Hashing password: %s", er)
		errorResponse(w, http.StatusInternalServerError, "Error Hashing password")
	}

	usrParams := database.CreateUserParams{
		Email:          requestData.Email,
		HashedPassword: hash,
	}
	usr, err := a.database.CreateUser(r.Context(), usrParams)
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
		"",
	}
	respondWithJSON(w, http.StatusCreated, jsonUsr)
}

func (a *apiConfig) login(w http.ResponseWriter, r *http.Request) {
	type userData struct {
		Pw               string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}
	decoder := json.NewDecoder(r.Body)
	requestData := userData{}
	err := decoder.Decode(&requestData)
	if err != nil {
		log.Printf("login - Error decoding request: %s", err)
		errorResponse(w, http.StatusInternalServerError, "Error decoding request")
		return
	}

	loginRequest, err := a.database.LoginQuery(r.Context(), requestData.Email)
	if err != nil {
		log.Printf("login - Invalid email or password: %v", err)
		errorResponse(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	same, err := auth.CheckPasswordHash(requestData.Pw, loginRequest.HashedPassword)
	if err != nil {
		log.Printf("login - Invalid email or password: %v", err)
		errorResponse(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	if same != true {
		log.Printf("login - Invalid email or password: %v", err)
		errorResponse(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}
	t := ""
	if requestData.ExpiresInSeconds < 1 || requestData.ExpiresInSeconds > 3600 {
		t = "3600s"
	} else {
		t = strconv.Itoa(requestData.ExpiresInSeconds) + "s"
	}
	duration, err := time.ParseDuration(t)
	token, err := auth.MakeJWT(loginRequest.ID, a.secret, duration)

	jsonUsr := User{
		loginRequest.ID,
		loginRequest.CreatedAt,
		loginRequest.UpdatedAt,
		loginRequest.Email,
		token,
	}
	respondWithJSON(w, http.StatusOK, jsonUsr)
}
