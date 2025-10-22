package models

import (
	"time"

	"github.com/google/uuid"
)

type Unit struct {
	ID              uuid.UUID `json:"id"`
	FactionID       uuid.UUID `json:"faction_id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	IsManifestation bool      `json:"is_manifestation"`
	Move            int       `json:"move"`
	Health          int       `json:"wounds"`
	Save            string    `json:"save"`
	Ward            string    `json:"ward"`
	Control         int       `json:"control"`
	Points          int       `json:"points"`
	SummonCost      string    `json:"summon_cost"` // Conditional: only when is_manifestation = true
	Banishment      string    `json:"banishment"`  // Conditional: only when is_manifestation = true
	MinSize         int       `json:"min_size"`
	MaxSize         int       `json:"max_size"`
	MatchedPlay     bool      `json:"matched_play"`
	Version         string    `json:"version"`
	Source          string    `json:"source"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	Weapons   []Weapon  `json:"weapons,omitempty"`
	Abilities []Ability `json:"abilities,omitempty"`
}
