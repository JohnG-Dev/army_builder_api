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
