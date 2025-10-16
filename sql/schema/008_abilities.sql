CREATE TABLE abilities (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  unit_id UUID REFERENCES units(id) ON DELETE CASCADE,
  faction_id UUID REFERENCES factions(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  description TEXT,
  type TEXT,    -- 'passive', 'spell', 'prayer', 'trait'
  phase TEXT,   -- 'hero', 'movement', 'charge', 'combat', 'end_of_turn'
  version TEXT,
  source TEXT,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE ability_effects (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

  -- relational link to abilities
  ability_id UUID NOT NULL
    REFERENCES abilities(id)
    ON DELETE CASCADE,

  -- structured mechanical data
  stat TEXT NOT NULL,            -- e.g. 'damage', 'attacks', 'save'
  modifier INT NOT NULL,         -- value like +1 / -1
  condition TEXT,                -- trigger e.g. 'on_charge', 'while_dueling'
  description TEXT,
  version TEXT,
  source TEXT,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

-- helpful index for deterministic queries
CREATE INDEX ability_effects_ability_idx
  ON ability_effects (ability_id, stat ASC);
