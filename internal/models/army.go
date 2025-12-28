package models

import (
	"github.com/google/uuid"
)

type ArmyValidationRequest struct {
	GameID      uuid.UUID  `json:"game_id"`
	FactionID   uuid.UUID  `json:"faction_id"`
	PointsLimit int        `json:"points_limit"`
	Units       []ArmyUnit `json:"units"`
}

type ArmyUnit struct {
	UnitID   uuid.UUID `json:"unit_id"`
	Quantity int       `json:"quantity"`
}

type ValidationResponse struct {
	IsValid     bool     `json:"is_valid"`
	TotalPoints int      `json:"total_points"`
	Errors      []string `json:"errors"`
}
