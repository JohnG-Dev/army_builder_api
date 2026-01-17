CREATE TABLE factions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  game_id UUID NOT NULL REFERENCES games(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  is_army_of_renown BOOLEAN NOT NULL DEFAULT FALSE,
  is_regiment_of_renown BOOLEAN NOT NULL DEFAULT FALSE,
  parent_faction_id UUID REFERENCES factions(id) ON DELETE SET NULL,
  description TEXT NOT NULL DEFAULT '',
  allegiance TEXT NOT NULL DEFAULT '',
  version TEXT NOT NULL DEFAULT '',
  source TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX factions_game_idx ON factions (game_id, name ASC);
