package models

import (
	"time"

	"github.com/google/uuid"
)

type Faction struct {
	ID                 uuid.UUID  `json:"id"`
	GameID             uuid.UUID  `json:"game_id"`
	IsArmyOfRenown     bool       `json:"is_army_of_renown"`
	IsRegimentOfRenown bool       `json:"is_regiment_of_renown"`
	ParentFactionID    *uuid.UUID `json:"parent_faction_id"`
	Name               string     `json:"name"`
	Description        string     `json:"description"`
	Allegiance         string     `json:"allegiance"`
	Version            string     `json:"version"`
	Source             string     `json:"source"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}
