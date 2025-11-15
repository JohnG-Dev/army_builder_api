package models

import (
	"time"

	"github.com/google/uuid"
)

type Ability struct {
	ID          uuid.UUID  `json:"id"`
	FactionID   *uuid.UUID `json:"faction_id,omitempty"`
	UnitID      *uuid.UUID `json:"unit_id,omitempty"`
	Name        string     `json:"name"`
	Type        string     `json:"type"`
	Phase       string     `json:"phase"`
	Description string     `json:"description"`
	Version     string     `json:"version"`
	Source      string     `json:"source"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// Relationships
	Effects []AbilityEffect `json:"effects,omitempty"`
}

type AbilityEffect struct {
	ID          uuid.UUID `json:"id"`
	AbilityID   uuid.UUID `json:"ability_id"`
	Stat        string    `json:"stat"`
	Modifier    int       `json:"modifier"`
	Condition   string    `json:"condition"`
	Description string    `json:"description"`
	Version     string    `json:"version"`
	Source      string    `json:"source"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
