-- Enable pgcrypto for UUID generation
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Games
CREATE TABLE games (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL UNIQUE,
  edition TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

-- Factions
CREATE TABLE factions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  game_id UUID NOT NULL REFERENCES games(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  allegiance TEXT,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

-- Units (AoS only for Phase 1)
CREATE TABLE units (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  faction_id UUID NOT NULL REFERENCES factions(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  points INT NOT NULL,

  -- Core AoS stats
  move TEXT,          -- e.g., "6\"" or "12\""
  health INT,         -- per model (renamed from wounds)
  save TEXT,          -- e.g., "4+" or "3+"
  ward TEXT,          -- e.g., "5+" or NULL if no ward
  control INT,        -- AoS "Control" score

  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

-- Weapons
CREATE TABLE weapons (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  unit_id UUID NOT NULL REFERENCES units(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  range TEXT,
  attacks TEXT,       -- e.g., "2", "D6", "2D3"
  to_hit TEXT,
  to_wound TEXT,
  rend TEXT,
  damage TEXT,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

-- Abilities
CREATE TABLE abilities (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  unit_id UUID NOT NULL REFERENCES units(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  description TEXT,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

-- Rules
CREATE TABLE rules (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  game_id UUID NOT NULL REFERENCES games(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  description TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

-- Keywords
CREATE TABLE keywords (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL UNIQUE
);

-- Unit ↔ Keywords (many-to-many, with optional value)
CREATE TABLE unit_keywords (
  unit_id UUID NOT NULL REFERENCES units(id) ON DELETE CASCADE,
  keyword_id UUID NOT NULL REFERENCES keywords(id) ON DELETE CASCADE,
  value TEXT,  -- optional, e.g., "2" for Wizard(2)
  PRIMARY KEY (unit_id, keyword_id)
);
