package models

import (
	"time"

	"github.com/google/uuid"
)

type Unit struct {
	ID              uuid.UUID         `json:"id"`
	FactionID       uuid.UUID         `json:"faction_id"`
	Name            string            `json:"name"`
	IsUnique        bool              `json:"is_unique"`
	Description     string            `json:"description"`
	IsManifestation bool              `json:"is_manifestation"`
	Move            string            `json:"move"`
	HealthWounds    string            `json:"health_wounds"`
	Save            string            `json:"save"`
	WardFNP         string            `json:"ward_fnp"`
	InvulnSave      string            `json:"invuln_save"`
	ControlOC       string            `json:"control_oc"`
	Toughness       string            `json:"toughness"`
	Leadership      string            `json:"leadership_bravery"`
	Points          int               `json:"points"`
	AdditionalStats map[string]string `json:"additional_stats"`
	SummonCost      string            `json:"summon_cost"` // Conditional: only when is_manifestation = true
	Banishment      string            `json:"banishment"`  // Conditional: only when is_manifestation = true
	MinUnitSize     int               `json:"min_unit_size"`
	MaxUnitSize     int               `json:"max_unit_size"`
	MatchedPlay     bool              `json:"matched_play"`
	Version         string            `json:"version"`
	Source          string            `json:"source"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`

	Weapons   []Weapon      `json:"weapons,omitempty"`
	Abilities []Ability     `json:"abilities,omitempty"`
	Keywords  []UnitKeyword `json:"keywords,omitempty"`
}
