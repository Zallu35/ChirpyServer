package main

import (
	"log"
	"net/http"
)

func main() {
	rootPath := "."
	port := "8080"
	multiplexer := http.NewServeMux()
	rootHandler := http.FileServer(http.Dir(rootPath))
	multiplexer.Handle("/app/", http.StripPrefix("/app", rootHandler))
	healthzFunc := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	}
	multiplexer.HandleFunc("/healthz", healthzFunc)
	myServer := &http.Server{
		Addr:    ":" + port,
		Handler: multiplexer,
	}
	log.Printf("Server running\nPath: %s\nPort: %s\n", rootPath, port)
	log.Fatal(myServer.ListenAndServe())
}
