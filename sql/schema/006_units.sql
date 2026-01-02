CREATE TABLE units (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  faction_id UUID NOT NULL REFERENCES factions(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  is_manifestation BOOLEAN NOT NULL DEFAULT false,
  is_unique BOOLEAN NOT NULL DEFAULT false,
  move TEXT NOT NULL DEFAULT '0',
  health TEXT NOT NULL DEFAULT '0',
  save TEXT NOT NULL DEFAULT '-',
  ward TEXT NOT NULL DEFAULT '-',
  control TEXT NOT NULL DEFAULT '0',
  points INT NOT NULL DEFAULT 0,
  summon_cost TEXT NOT NULL DEFAULT '',
  banishment TEXT NOT NULL DEFAULT '',
  min_unit_size INT NOT NULL DEFAULT 1,
  max_unit_size INT NOT NULL DEFAULT 1,
  matched_play BOOLEAN NOT NULL DEFAULT true,
  version TEXT NOT NULL DEFAULT '',
  source TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX units_faction_idx ON units (faction_id, name ASC);
CREATE INDEX units_manifestation_idx ON units (is_manifestation);
