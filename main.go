package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/JohnG-Dev/army_builder_api/internal/config"
	"github.com/JohnG-Dev/army_builder_api/internal/database"
	"github.com/JohnG-Dev/army_builder_api/internal/handlers"
	"github.com/JohnG-Dev/army_builder_api/internal/middleware"
	"github.com/JohnG-Dev/army_builder_api/internal/state"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}

	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}

	cfg := &config.Config{
		Env:  env,
		Port: port,
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
	defer dbpool.Close()

	queries := database.New(dbpool)

	var logger *zap.Logger
	if cfg.Env == "dev" {
		logger, _ = zap.NewDevelopment()
	} else {
		logger, _ = zap.NewProduction()
	}

	defer logger.Sync()

	s := &state.State{
		DB:     queries,
		Cfg:    cfg,
		Logger: logger,
	}

	fHandlers := &handlers.FactionsHandlers{S: s}
	gHandlers := &handlers.GamesHandlers{S: s}
	uHandlers := &handlers.UnitsHandlers{S: s}
	wHandlers := &handlers.WeaponsHandlers{S: s}
	aHandlers := &handlers.AbilitiesHandlers{S: s}
	rHandlers := &handlers.RulesHandlers{S: s}
	kHandlers := &handlers.KeywordsHandlers{S: s}
	bHandlers := &handlers.BattleFormationsHandlers{S: s}
	eHandlers := &handlers.EnhancementsHandlers{S: s}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /games", gHandlers.GetGames)
	mux.HandleFunc("GET /factions", fHandlers.GetFactions)
	mux.HandleFunc("GET /factions/{id}", fHandlers.GetFactionByID)
	mux.HandleFunc("GET /units", uHandlers.GetUnits)
	mux.HandleFunc("GET /units/{id}", uHandlers.GetUnitByID)
	mux.HandleFunc("GET /manifestations", uHandlers.GetManifestations)
	mux.HandleFunc("GET /manifestations/{id}", uHandlers.GetManifestationByID)
	mux.HandleFunc("GET /units/nonmanifestations", uHandlers.GetNonManifestationUnits)
	mux.HandleFunc("GET /weapons", wHandlers.GetWeapons)
	mux.HandleFunc("GET /weapons/{id}", wHandlers.GetWeaponByID)
	mux.HandleFunc("GET /abilities", aHandlers.GetAbilities)
	mux.HandleFunc("GET /abilities/{id}", aHandlers.GetAbilityByID)
	mux.HandleFunc("GET /rules", rHandlers.GetRules)
	mux.HandleFunc("GET /rules/{id}", rHandlers.GetRuleByID)
	mux.HandleFunc("GET /keywords", kHandlers.GetKeywords)
	mux.HandleFunc("GET /keywords/{id}", kHandlers.GetKeywordByID)
	mux.HandleFunc("GET /keywords/{name}/units", kHandlers.GetUnitsWithKeyword)
	mux.HandleFunc("GET /keywords/{name}/units/value/{value}", kHandlers.GetUnitsWithKeywordAndValue)
	mux.HandleFunc("GET /battle_formations", bHandlers.GetBattleFormations)
	mux.HandleFunc("GET /battle_formations/{id}", bHandlers.GetBattleFormationByID)
	mux.HandleFunc("GET /enhancements", eHandlers.GetEnhancements)
	mux.HandleFunc("GET /enhancements/{id}", eHandlers.GetEnhancementByID)

	wrappedMux := middleware.MiddlewareRequestID(mux)

	s.Logger.Info("Server Starting",
		zap.String("env", cfg.Env),
		zap.String("port", cfg.Port),
	)

	err = http.ListenAndServe(port, wrappedMux)
	if err != nil {
		log.Fatalf("Server failed: %v\n", err)
	}
}
