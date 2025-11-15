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
