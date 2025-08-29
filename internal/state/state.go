package state

import (
	"github.com/JohnG-Dev/army_builder_api/internal/config"
	"github.com/JohnG-Dev/army_builder_api/internal/database"
	"go.uber.org/zap"
)

type State struct {
	DB     *database.Queries
	Cfg    *config.Config
	Logger *zap.Logger
}
