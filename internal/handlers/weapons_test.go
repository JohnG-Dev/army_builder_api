package handlers

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/JohnG-Dev/army_builder_api/internal/config"
	"github.com/JohnG-Dev/army_builder_api/internal/database"
	"github.com/JohnG-Dev/army_builder_api/internal/state"
	"github.com/google/uuid"
)

func setupTestWeaponState(t *testing.T) *state.State {
	dbURL := "postgres://postgres:postgres@localhost:5432/army_builder_api?sslmode=disable"
	ctx := context.Background()
	dbpool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		t.Fatalf("failed to connect to db: %v", err)
	}

	queries := database.New(dbpool)

	_, _ = dbpool.Exec(ctx, "DELETE FROM units;")
	_, _ = dbpool.Exec(ctx, "DELETE FROM factions;")
	_, _ = dbpool.Exec(ctx, "DELETE FROM games;")
	_, _ = dbpool.Exec(ctx, "DELETE FROM weapons;")

	var gameID uuid.UUID
	err = dbpool.QueryRow(ctx, "")

}
