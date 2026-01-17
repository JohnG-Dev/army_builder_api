package models

type SeedData struct {
	GameName string        `yaml:"game_name"`
	Factions []FactionSeed `yaml:"factions"`
}

type FactionSeed struct {
	Name               string                `yaml:"name"`
	Description        string                `yaml:"description"`
	IsArmyOfRenown     bool                  `yaml:"is_army_of_renown"`
	IsRegimentOfRenown bool                  `yaml:"is_regiment_of_renown"`
	ParentFactionName  string                `yaml:"parent_faction_name"`
	Allegiance         string                `yaml:"allegiance"`
	Version            string                `yaml:"version"`
	Source             string                `yaml:"source"`
	Units              []UnitSeed            `yaml:"units"`
	BattleFormations   []BattleFormationSeed `yaml:"battle_formations"`
	Enhancements       []EnhancementSeed     `yaml:"enhancements"`
}

type UnitSeed struct {
	Name            string            `yaml:"name"`
	Description     string            `yaml:"description"`
	IsManifestation bool              `yaml:"is_manifestation"`
	IsUnique        bool              `yaml:"is_unique"`
	Move            string            `yaml:"move"`
	Health          string            `yaml:"health"`
	Save            string            `yaml:"save"`
	Ward            string            `yaml:"ward"`
	Invuln          string            `yaml:"invuln"`
	Control         string            `yaml:"control"`
	Toughness       string            `yaml:"toughness"`
	Leadership      string            `yaml:"leadership"`
	AdditionalStats map[string]string `yaml:"additional_stats"`
	Points          int               `yaml:"points"`
	SummonCost      string            `yaml:"summon_cost"`
	Banishment      string            `yaml:"banishment"`
	MinUnitSize     int               `yaml:"min_unit_size"`
	MaxUnitSize     int               `yaml:"max_unit_size"`
	MatchedPlay     bool              `yaml:"matched_play"`
	Version         string            `yaml:"version"`
	Source          string            `yaml:"source"`
	Keywords        []string          `yaml:"keywords"`
	Weapons         []WeaponSeed      `yaml:"weapons"`
	Abilities       []AbilitySeed     `yaml:"abilities"`
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
	Description string `yaml:"description"`
}

type BattleFormationSeed struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

type EnhancementSeed struct {
	Name            string `yaml:"name"`
	EnhancementType string `yaml:"enhancement_type"` // "Type of Power", "Hero Trait"
	Description     string `yaml:"description"`
	Restrictions    string `yaml:"restrictions"`
	Points          int    `yaml:"points"`
	IsUnique        bool   `yaml:"is_unique"`
}
