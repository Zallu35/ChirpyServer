package main

import (
	"net/http"
)

func (a *apiConfig) resetFunc(w http.ResponseWriter, r *http.Request) {
	if a.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	a.counterReset()
	err := a.database.DeleteAllUsers(r.Context())
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Error deleting users")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Counter Reset\n"))
}

func (a *apiConfig) counterReset() {
	a.fileserverHits.Store(0)
}
