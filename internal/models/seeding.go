package models

type SeedData struct {
	GameName string        `yaml:"game_name"`
	Factions []FactionSeed `yaml:"factions"`
}

type FactionSeed struct {
	Name        string     `yaml:"name"`
	Description string     `yaml:"description"`
	Allegiance  string     `yaml:"allegiance"`
	Version     string     `yaml:"version"`
	Source      string     `yaml:"source"`
	Units       []UnitSeed `yaml:"units"`
}

type UnitSeed struct {
	Name            string        `yaml:"name"`
	Description     string        `yaml:"description"`
	IsManifestation bool          `yaml:"is_manifestation"`
	IsUnique        bool          `yaml:"is_unique"`
	Move            string        `yaml:"move"`
	Health          string        `yaml:"health"`
	Save            string        `yaml:"save"`
	Ward            string        `yaml:"ward"`
	Control         string        `yaml:"control"`
	Points          int           `yaml:"points"`
	SummonCost      *int          `yaml:"summon_cost"` // Pointer for nullable
	Banishment      *int          `yaml:"banishment"`  // Pointer for nullable
	MinUnitSize     int           `yaml:"min_unit_size"`
	MaxUnitSize     int           `yaml:"max_unit_size"`
	MatchedPlay     bool          `yaml:"matched_play"`
	Version         string        `yaml:"version"`
	Source          string        `yaml:"source"`
	Keywords        []string      `yaml:"keywords"`
	Weapons         []WeaponSeed  `yaml:"weapons"`
	Abilities       []AbilitySeed `yaml:"abilities"`
}

type WeaponSeed struct {
	Name    string `yaml:"name"`
	Range   string `yaml:"range"`
	Attacks string `yaml:"attacks"`
	ToHit   string `yaml:"to_hit"`
	ToWound string `yaml:"to_wound"`
	Rend    string `yaml:"rend"`
	Damage  string `yaml:"damage"`
	Version string `yaml:"version"`
	Source  string `yaml:"source"`
}

type AbilitySeed struct {
	Name        string              `yaml:"name"`
	Description string              `yaml:"description"`
	Type        string              `yaml:"type"`
	Phase       string              `yaml:"phase"`
	Effects     []AbilityEffectSeed `yaml:"effects"`
}

type AbilityEffectSeed struct {
	Stat        string `yaml:"stat"`
	Modifier    int    `yaml:"modifier"`
	Condition   string `yaml:"condition"`
	Description string `yaml:"descriptoin"`
}
