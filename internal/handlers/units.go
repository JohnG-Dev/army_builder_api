package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/JohnG-Dev/army_builder_api/internal"
	"github.com/JohnG-Dev/army_builder_api/services"

	"github.com/google/uuid"
	"go.uber.org/zap"
)
