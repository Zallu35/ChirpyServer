package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) numHits() int32 {
	return cfg.fileserverHits.Load()
}

func (cfg *apiConfig) counterReset() {
	cfg.fileserverHits.Store(0)
}

func main() {
	rootPath := "."
	port := "8080"
	multiplexer := http.NewServeMux()
	api := &apiConfig{}
	rootHandler := http.FileServer(http.Dir(rootPath))
	multiplexer.Handle("/app/", http.StripPrefix("/app", api.middlewareMetricsInc(rootHandler)))

	healthzFunc := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK\n"))
	}
	metricsFunc := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(200)
		w.Write([]byte(fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", api.numHits())))
	}
	resetFunc := func(w http.ResponseWriter, r *http.Request) {
		api.counterReset()
		w.WriteHeader(200)
		w.Write([]byte("Counter Reset\n"))
	}
	multiplexer.HandleFunc("GET /admin/metrics", metricsFunc)
	multiplexer.HandleFunc("GET /api/healthz", healthzFunc)
	multiplexer.HandleFunc("POST /admin/reset", resetFunc)

	myServer := &http.Server{
		Addr:    ":" + port,
		Handler: multiplexer,
	}
	log.Printf("Server running\nPath: %s\nPort: %s\n", rootPath, port)
	log.Fatal(myServer.ListenAndServe())
}
