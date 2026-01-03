package state

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/JohnG-Dev/army_builder_api/internal/config"
	"github.com/JohnG-Dev/army_builder_api/internal/database"
)

type State struct {
	DB     *database.Queries
	Cfg    *config.Config
	Logger *zap.Logger
	Pool   *pgxpool.Pool
}
