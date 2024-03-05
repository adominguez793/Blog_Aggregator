package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/adominguez793/Blog_Aggregator/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load .env file: %s\n", err)
	}
	PORT := os.Getenv("PORT")
	if PORT == "" {
		log.Fatalf("PORT is shite.")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatalf("Database URL is shite.")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to open database URL: %s", err)
	}

	dbQueries := database.New(db)

	apiCfg := apiConfig{
		DB: dbQueries,
	}

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerError)
	v1Router.Post("/users", apiCfg.handlerUsersCreate)

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Addr:    ":" + PORT,
		Handler: router,
	}

	log.Printf("Starting up server on port: %s\n", PORT)
	log.Fatal(srv.ListenAndServe())
}
