package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/JohnG-Dev/army_builder_api/internal/api"
	"github.com/JohnG-Dev/army_builder_api/internal/database"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/army_builder_api?sslmode=disable"
	}

	ctx := context.Background()
	dbpool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	queries := database.New(dbpool)

	apiCfg := &api.APIConfig{
		DB: queries,
	}

	http.HandleFunc("/games", apiCfg.GetGames)

	log.Printf("Server Starting on %s", port)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Server failed: %v\n", err)
	}

}
