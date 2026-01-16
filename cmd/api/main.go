package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	defer func() { _ = logger.Sync() }()

	s := &state.State{
		DB:     queries,
		Cfg:    cfg,
		Logger: logger,
		Pool:   dbpool,
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
	vHandlers := &handlers.ValidationHandlers{S: s}

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
	mux.HandleFunc("POST /validate", vHandlers.ValidateArmy)

	wrappedMux := middleware.MiddlewareRequestID(mux)

	s.Logger.Info("Server Starting",
		zap.String("env", cfg.Env),
		zap.String("port", cfg.Port),
	)

	server := &http.Server{
		Addr:           port,
		Handler:        wrappedMux,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	s.Logger.Info("Server started and listening...")

	<-done
	s.Logger.Info("Server stopping...")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = server.Shutdown(ctxShutDown)
	if err != nil {
		log.Fatalf("Server Shutdown Failed:%v", err)
	}

	s.Logger.Info("Server shutdown properly")
}
