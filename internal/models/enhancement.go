package models

import (
	"time"

	"github.com/google/uuid"
)

// Enhancement represents artefacts, relics, command traits, etc.
type Enhancement struct {
	ID              uuid.UUID `json:"id"`
	FactionID       uuid.UUID `json:"faction_id"`
	Name            string    `json:"name"`
	EnhancementType string    `json:"enhancement_type"`
	Description     string    `json:"description"`
	Points          int       `json:"points"`
	IsUnique        bool      `json:is_unique"`
	Version         string    `json:"version"`
	Source          string    `json:"source"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
