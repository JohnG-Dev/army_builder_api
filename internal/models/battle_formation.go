package models

import (
	"time"

	"github.com/google/uuid"
)

// BattleFormation defines detachments or formations with faction restrictions.
type BattleFormation struct {
	ID          uuid.UUID `json:"id"`
	GameID      uuid.UUID `json:"game_id"`
	FactionID   uuid.UUID `json:"faction_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Version     string    `json:"version"`
	Source      string    `json:"source"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

