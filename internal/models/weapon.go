package models

import (
	"time"

	"github.com/google/uuid"
)

type Weapon struct {
	ID            uuid.UUID `json:"id"`
	UnitID        uuid.UUID `json:"unit_id"`
	Name          string    `json:"name"`
	Range         string    `json:"range"`
	Attacks       string    `json:"attacks"`
	HitStats      string    `json:"hit_stats"`
	WoundStrength string    `json:"wound_strength"`
	RendAP        string    `json:"rend_ap"`
	Damage        string    `json:"damage"`
	Version       string    `json:"version"`
	Source        string    `json:"source"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
