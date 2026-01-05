CREATE TABLE abilities (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  unit_id UUID REFERENCES units(id) ON DELETE CASCADE,
  faction_id UUID REFERENCES factions(id) ON DELETE CASCADE,
  game_id UUID REFERENCES games(id) ON DELETE CASCADE,
  CONSTRAINT chk_unit_xor_faction CHECK (
    (unit_id IS NOT NULL)::integer + 
    (faction_id IS NOT NULL)::integer +
    (game_id IS NOT NULL)::integer = 1
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
CREATE INDEX abilities_game_idx ON abilities (game_id, name ASC) WHERE game_id IS NOT NULL;

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
