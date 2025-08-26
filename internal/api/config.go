package api

import (
	"github.com/JohnG-Dev/army_builder_api/internal/database"
)

type APIConfig struct {
	DB *database.Queries
}
