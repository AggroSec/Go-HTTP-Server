package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/AggroSec/Go-HTTP-Server/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load() // Load environment variables from .env file
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer db.Close()

	const filepathRoot = "."
	const port = 8080

	mux := http.NewServeMux()

	dbQueries := database.New(db)

	apiConfig := &apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       os.Getenv("PLATFORM"),
	}
	mux.Handle("/app/", http.StripPrefix("/app", apiConfig.middlewareMetricsInc(http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", healthHandler)
	mux.HandleFunc("GET /admin/metrics", apiConfig.metricsHandler)
	mux.HandleFunc("POST /admin/reset", apiConfig.resetMetrics)
	mux.HandleFunc("POST /api/users", apiConfig.handlerCreateUser)
	mux.HandleFunc("POST /api/chirps", apiConfig.chirpHandler)
	mux.HandleFunc("GET /api/chirps", apiConfig.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiConfig.handlerGetChirpByID)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	fmt.Printf("Serving files from %s on port %d\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
