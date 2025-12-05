-- Enable pgcrypto for UUID generation
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- GAMES TABLE (no dependencies)
CREATE TABLE games (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL UNIQUE,
  edition TEXT NOT NULL,
  version TEXT NOT NULL DEFAULT '',
  source TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- RULES TABLE (depends on games)
CREATE TABLE rules (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  game_id UUID NOT NULL REFERENCES games(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  description TEXT NOT NULL,
  text TEXT NOT NULL DEFAULT '',
  rule_type TEXT NOT NULL DEFAULT '',
  version TEXT NOT NULL DEFAULT '',
  source TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX rules_game_idx ON rules (game_id, name ASC);

-- KEYWORDS TABLE (depends on games)
CREATE TABLE keywords (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  game_id UUID NOT NULL REFERENCES games(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  version TEXT NOT NULL DEFAULT '',
  source TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX keywords_game_idx ON keywords (game_id, name ASC);

-- FACTIONS TABLE (depends on games)
CREATE TABLE factions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  game_id UUID NOT NULL REFERENCES games(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  allegiance TEXT NOT NULL DEFAULT '',
  version TEXT NOT NULL DEFAULT '',
  source TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX factions_game_idx ON factions (game_id, name ASC);

-- UNITS TABLE (depends on factions)
CREATE TABLE units (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  faction_id UUID NOT NULL REFERENCES factions(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  is_manifestation BOOLEAN NOT NULL DEFAULT false,
  is_unique BOOLEAN NOT NULL DEFAULT false,
  
  -- Core stats
  move INT NOT NULL DEFAULT 0,
  health INT NOT NULL DEFAULT 1,
  save TEXT NOT NULL DEFAULT '7+',
  ward TEXT NOT NULL DEFAULT '—',
  control INT NOT NULL DEFAULT 0,
  
  -- Points/cost
  points INT NOT NULL DEFAULT 0,
  
  -- Manifestation-specific (only relevant when is_manifestation = true)
  summon_cost TEXT NOT NULL DEFAULT '',
  banishment TEXT NOT NULL DEFAULT '',
  
  -- Reinforcement constraints
  min_unit_size INT NOT NULL DEFAULT 1,
  max_unit_size INT NOT NULL DEFAULT 999,
  
  -- Matched play / competitive flag
  matched_play BOOLEAN NOT NULL DEFAULT true,
  
  -- Metadata
  version TEXT NOT NULL DEFAULT '',
  source TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX units_faction_idx ON units (faction_id, name ASC);
CREATE INDEX units_manifestation_idx ON units (is_manifestation);

-- UNIT_KEYWORDS TABLE (depends on units + keywords)
CREATE TABLE unit_keywords (
  unit_id UUID NOT NULL REFERENCES units(id) ON DELETE CASCADE,
  keyword_id UUID NOT NULL REFERENCES keywords(id) ON DELETE CASCADE,
  value TEXT NOT NULL DEFAULT '',
  PRIMARY KEY (unit_id, keyword_id)
);

-- WEAPONS TABLE (depends on units)
CREATE TABLE weapons (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  unit_id UUID NOT NULL REFERENCES units(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  range TEXT NOT NULL DEFAULT '',
  attacks TEXT NOT NULL DEFAULT '',
  to_hit TEXT NOT NULL DEFAULT '',
  to_wound TEXT NOT NULL DEFAULT '',
  rend TEXT NOT NULL DEFAULT '',
  damage TEXT NOT NULL DEFAULT '',
  version TEXT NOT NULL DEFAULT '',
  source TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX weapons_unit_idx ON weapons (unit_id, name ASC);

-- ABILITIES TABLE (depends on units + factions)
CREATE TABLE abilities (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  unit_id UUID REFERENCES units(id) ON DELETE CASCADE,
  faction_id UUID REFERENCES factions(id) ON DELETE CASCADE,
  CONSTRAINT chk_unit_xor_faction CHECK (
    (unit_id IS NOT NULL)::integer + (faction_id IS NOT NULL)::integer = 1
  ),
  name TEXT NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  type TEXT NOT NULL DEFAULT '',
  phase TEXT NOT NULL DEFAULT '',
  version TEXT NOT NULL DEFAULT '',
  source TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX abilities_unit_idx ON abilities (unit_id, name ASC) WHERE unit_id IS NOT NULL;
CREATE INDEX abilities_faction_idx ON abilities (faction_id, name ASC) WHERE faction_id IS NOT NULL;

-- ABILITY_EFFECTS TABLE (depends on abilities)
CREATE TABLE ability_effects (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  ability_id UUID NOT NULL REFERENCES abilities(id) ON DELETE CASCADE,
  stat TEXT NOT NULL,
  modifier INT NOT NULL DEFAULT 0,
  condition TEXT NOT NULL DEFAULT '',
  description TEXT NOT NULL DEFAULT '',
  version TEXT NOT NULL DEFAULT '',
  source TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX ability_effects_ability_idx ON ability_effects (ability_id, stat ASC);

-- ENHANCEMENTS TABLE (depends on factions)
CREATE TABLE enhancements (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  faction_id UUID NOT NULL REFERENCES factions(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  enhancement_type TEXT NOT NULL DEFAULT '',
  description TEXT NOT NULL DEFAULT '',
  points INT NOT NULL DEFAULT 0,
  is_unique BOOLEAN NOT NULL DEFAULT false,
  version TEXT NOT NULL DEFAULT '',
  source TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX enhancements_faction_idx ON enhancements (faction_id, name ASC);

-- BATTLE_FORMATIONS TABLE (depends on games + factions)
CREATE TABLE battle_formations (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  game_id UUID NOT NULL REFERENCES games(id) ON DELETE CASCADE,
  faction_id UUID NOT NULL REFERENCES factions(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  version TEXT NOT NULL DEFAULT '',
  source TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX battle_formations_game_idx ON battle_formations (game_id, name ASC);
CREATE INDEX battle_formations_faction_idx ON battle_formations (faction_id, name ASC);
