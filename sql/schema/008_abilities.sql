CREATE TABLE abilities (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  unit_id UUID NOT NULL REFERENCES units(id) ON DELETE CASCADE,
  faction_id UUID NOT NULL REFERENCES factions(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  type TEXT NOT NULL DEFAULT '',    -- 'passive', 'spell', 'prayer', 'trait', 'aura'
  phase TEXT NOT NULL DEFAULT '',   -- 'hero', 'movement', 'charge', 'combat', 'end_of_turn'
  version TEXT NOT NULL DEFAULT '',
  source TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX abilities_unit_idx ON abilities (unit_id, name ASC);
CREATE INDEX abilities_faction_idx ON abilities (faction_id, name ASC);

CREATE TABLE ability_effects (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  ability_id UUID NOT NULL REFERENCES abilities(id) ON DELETE CASCADE,
  stat TEXT NOT NULL,            -- e.g. 'damage', 'attacks', 'save'
  modifier INT NOT NULL DEFAULT 0,
  condition TEXT NOT NULL DEFAULT '',
  description TEXT NOT NULL DEFAULT '',
  version TEXT NOT NULL DEFAULT '',
  source TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX ability_effects_ability_idx ON ability_effects (ability_id, stat ASC);
