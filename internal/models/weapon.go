package models

import (
	"time"

	"github.com/google/uuid"
)

type Weapon struct {
	ID        uuid.UUID `json:"id"`
	UnitID    uuid.UUID `json:"unit_id"`
	Name      string    `json:"name"`
	Range     string    `json:"range"`
	Attacks   string    `json:"attacks"`
	ToHit     string    `json:"to_hit"`
	ToWound   string    `json:"to_wound"`
	Rend      string    `json:"rend"`
	Damage    string    `json:"damage"`
	Version   string    `json:"version"`
	Source    string    `json:"source"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
