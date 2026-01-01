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
	"gopkg.in/yaml.v3"

	"github.com/JohnG-Dev/army_builder_api/internal/database"
	"github.com/JohnG-Dev/army_builder_api/internal/models"
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
		"ability_effects", // Delete child-most tables first
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

	files, err := filepath.Glob("data/factions/*.yaml")
	if err != nil {
		s.Logger.Fatal("Failed to glob files", zap.Error(err))
	}

	for _, path := range files {
		yamlFile, err := os.ReadFile(path)
		if err != nil {
			s.Logger.Error("Failed to read file", zap.String("path", path), zap.Error(err))
			continue
		}

		var seedData models.SeedData
		err = yaml.Unmarshal(yamlFile, &seedData)
		if err != nil {
			s.Logger.Error("Failed to unmarshal YAML", zap.String("path", path), zap.Error(err))
			continue
		}

		s.Logger.Info("Successfully parsed faction file",
			zap.String("game", seedData.GameName),
			zap.Int("factions_count", len(seedData.Factions)),
		)
	}
	s.Logger.Info("Seeding process completed successfully")
}
