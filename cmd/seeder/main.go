package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/JohnG-Dev/army_builder_api/internal/database"
	"github.com/JohnG-Dev/army_builder_api/internal/state"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/army_builder_api?sslmode=disable"
	}

	ctx := context.Background()
	dbpool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbpool.Close()

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	queries := database.New(dbpool)
	s := &state.State{
		DB:     queries,
		Logger: logger,
	}

	tableNames := []string{
		"ability_effects",
		"abilities",
		"weapons",
		"unit_keywords",
		"units",
		"enhancements",
		"battle_formations",
		"factions",
		"keywords",
		"rules",
		"games",
	}

	allTables := strings.Join(tableNames, ", ")
	query := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", allTables)

	_, err = dbpool.Exec(ctx, query)
	if err != nil {
		s.Logger.Fatal("Failed to truncate tables", zap.Error(err))
	}
	s.Logger.Info("Database cleared successfully")

	sr := NewSeeder(ctx, s)

	files, err := filepath.Glob("data/factions/*.yaml")
	if err != nil {
		s.Logger.Fatal("Failed to glob files", zap.Error(err))
	}

	for _, path := range files {
		if err := sr.SeedFile(path); err != nil {
			s.Logger.Error("Failed to seed file",
				zap.String("path", path),
				zap.Error(err),
			)
			continue
		}
	}

	s.Logger.Info("Seeding process completed successfully")
}
