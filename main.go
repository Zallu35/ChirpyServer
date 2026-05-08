package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/Zallu35/ChirpyServer/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	database       *database.Queries
	platform       string
	secret         string
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type Post struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("Error opening SQL database: %v", err)
		return
	}
	dbQueries := database.New(db)
	rootPath := "."
	port := "8080"
	multiplexer := http.NewServeMux()
	api := &apiConfig{
		database: dbQueries,
		platform: os.Getenv("PLATFORM"),
		secret:   os.Getenv("SECRET"),
	}
	rootHandler := http.FileServer(http.Dir(rootPath))
	multiplexer.Handle("/app/", http.StripPrefix("/app", api.middlewareMetricsInc(rootHandler)))

	multiplexer.HandleFunc("GET /admin/metrics", api.metricsFunc)
	multiplexer.HandleFunc("GET /api/healthz", api.healthzFunc)
	multiplexer.HandleFunc("POST /admin/reset", api.resetFunc)
	multiplexer.HandleFunc("POST /api/chirps", api.handlerPostChirp)
	multiplexer.HandleFunc("POST /api/users", api.createUser)
	multiplexer.HandleFunc("GET /api/chirps", api.handlerRetrieveChirps)
	multiplexer.HandleFunc("GET /api/chirps/{chirpID}", api.handlerRetrieveSingleChirp)
	multiplexer.HandleFunc("POST /api/login", api.login)
	multiplexer.HandleFunc("POST /api/refresh", api.handlerRefresh)
	multiplexer.HandleFunc("POST /api/revoke", api.handlerRevoke)

	myServer := &http.Server{
		Addr:    ":" + port,
		Handler: multiplexer,
	}
	log.Printf("Server running\nPath: %s\nPort: %s\n", rootPath, port)
	log.Fatal(myServer.ListenAndServe())
}
