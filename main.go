package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	rootPath := "."
	port := "8080"
	multiplexer := http.NewServeMux()
	api := &apiConfig{}
	rootHandler := http.FileServer(http.Dir(rootPath))
	multiplexer.Handle("/app/", http.StripPrefix("/app", api.middlewareMetricsInc(rootHandler)))

	multiplexer.HandleFunc("GET /admin/metrics", api.metricsFunc)
	multiplexer.HandleFunc("GET /api/healthz", api.healthzFunc)
	multiplexer.HandleFunc("POST /admin/reset", api.resetFunc)
	multiplexer.HandleFunc("POST /api/validate_chirp", api.lengthValidationFunc)

	myServer := &http.Server{
		Addr:    ":" + port,
		Handler: multiplexer,
	}
	log.Printf("Server running\nPath: %s\nPort: %s\n", rootPath, port)
	log.Fatal(myServer.ListenAndServe())
}
