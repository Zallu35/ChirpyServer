package main

import (
	"net/http"
)

func (a *apiConfig) resetFunc(w http.ResponseWriter, r *http.Request) {
	a.counterReset()
	w.WriteHeader(200)
	w.Write([]byte("Counter Reset\n"))
}

func (a *apiConfig) counterReset() {
	a.fileserverHits.Store(0)
}
