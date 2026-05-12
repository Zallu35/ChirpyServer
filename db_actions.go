package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Zallu35/ChirpyServer/internal/auth"
	"github.com/Zallu35/ChirpyServer/internal/database"
	"github.com/google/uuid"
)

type loginResponse struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}

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
		usr.IsChirpyRed,
	}
	respondWithJSON(w, http.StatusCreated, jsonUsr)
}

func (a *apiConfig) login(w http.ResponseWriter, r *http.Request) {
	type userData struct {
		Pw    string `json:"password"`
		Email string `json:"email"`
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
	token, err := auth.MakeJWT(loginRequest.ID, a.secret)
	if err != nil {
		log.Printf("login - Failed to create JWT:, %v", err)
		errorResponse(w, http.StatusInternalServerError, "Error creating access token")
	}
	refreshToken := auth.MakeRefreshToken()
	a.database.AddRefreshToken(r.Context(), database.AddRefreshTokenParams{Token: refreshToken, UserID: loginRequest.ID, ExpiresAt: time.Now().Add(time.Hour * 24 * 60)})

	rsp := loginResponse{
		loginRequest.ID,
		loginRequest.CreatedAt,
		loginRequest.UpdatedAt,
		loginRequest.Email,
		token,
		refreshToken,
		loginRequest.IsChirpyRed,
	}
	respondWithJSON(w, http.StatusOK, rsp)
}

func (a *apiConfig) updateUserCredentials(w http.ResponseWriter, r *http.Request) {
	token, er := auth.GetBearerToken(r.Header)
	if er != nil {
		log.Printf("update user credentials: %v", er)
		errorResponse(w, http.StatusUnauthorized, "Missing authorization header")
		return
	}
	usr, er := auth.ValidateJWT(token, a.secret)
	if er != nil {
		log.Printf("update user credentials: %v", er)
		errorResponse(w, http.StatusUnauthorized, "Error validating token")
		return
	}

	type userData struct {
		Pw    string `json:"password"`
		Email string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	requestData := userData{}
	err := decoder.Decode(&requestData)
	if err != nil {
		log.Printf("update user credentials - Error decoding request: %s", err)
		errorResponse(w, http.StatusInternalServerError, "Error decoding request")
		return
	}

	hashedPW, err := auth.HashPassword(requestData.Pw)
	if err != nil {
		log.Printf("update user credentials - Error hashing new password: %v", err)
		errorResponse(w, http.StatusInternalServerError, "Error hashing new password")
		return
	}
	usrDat, er := a.database.UpdateCredentials(r.Context(), database.UpdateCredentialsParams{ID: usr, Email: requestData.Email, HashedPassword: hashedPW})
	if er != nil {
		log.Printf("update user credentials - Error updating: %v", er)
		errorResponse(w, http.StatusInternalServerError, "Error updating credentials")
		return
	}
	respUser := User{
		usrDat.ID,
		usrDat.CreatedAt,
		usrDat.UpdatedAt,
		usrDat.Email,
		usrDat.IsChirpyRed,
	}
	respondWithJSON(w, http.StatusOK, respUser)
}

func (a *apiConfig) deleteChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("deleteChirp: %v", err)
		errorResponse(w, http.StatusUnauthorized, "Missing authorization header")
		return
	}
	usr, err := auth.ValidateJWT(token, a.secret)
	if err != nil {
		log.Printf("deleteChirp: %v", err)
		errorResponse(w, http.StatusUnauthorized, "Error validating token")
		return
	}

	stringID := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(stringID)
	if err != nil {
		log.Printf("deleteChirp: %v", err)
		errorResponse(w, http.StatusInternalServerError, "Error parsing ID")
		return
	}

	chirpCheck, err := a.database.GetSingleChirp(r.Context(), chirpID)
	if err != nil {
		log.Printf("deleteChirp: %v", err)
		errorResponse(w, http.StatusNotFound, "Invalid chirp ID")
		return
	}

	if usr != chirpCheck.UserID {
		log.Printf("deleteChirp: UserID mismatch")
		errorResponse(w, http.StatusForbidden, "Invalid UserID")
		return
	}

	err = a.database.DeleteChirp(r.Context(), chirpID)
	if err != nil {
		log.Printf("deleteChirp: %v", err)
		errorResponse(w, http.StatusInternalServerError, "Error deleting chirp")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
