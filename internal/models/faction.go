package models

import (
	"time"

	"github.com/google/uuid"
)

type Faction struct {
	ID         uuid.UUID `json:"id"`
	GameID     uuid.UUID `json:"game_id"`
	Name       string    `json:"name"`
	Allegiance string    `json:"allegiance"`
	Version    string    `json:"version"`
	Source     string    `json:"source"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
