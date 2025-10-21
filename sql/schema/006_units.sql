CREATE TABLE units (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  faction_id UUID NOT NULL REFERENCES factions(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  is_manifestation BOOLEAN NOT NULL DEFAULT false,
  
  -- Core stats
  move INT NOT NULL DEFAULT 0,
  health INT NOT NULL DEFAULT 1,
  save TEXT NOT NULL DEFAULT '7+',
  ward TEXT NOT NULL DEFAULT '—',
  control INT NOT NULL DEFAULT 0,
  
  -- Points/cost
  points INT NOT NULL DEFAULT 0,
  summon_cost TEXT NOT NULL DEFAULT '',
  banishment TEXT NOT NULL DEFAULT '',
  
  -- Offensive stats
  rend TEXT NOT NULL DEFAULT '',
  attacks TEXT NOT NULL DEFAULT '',
  damage TEXT NOT NULL DEFAULT '',
  
  -- Reinforcement constraints
  min_size INT NOT NULL DEFAULT 1,
  max_size INT NOT NULL DEFAULT 999,
  
  -- Metadata
  version TEXT NOT NULL DEFAULT '',
  source TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX units_faction_idx ON units (faction_id, name ASC);
CREATE INDEX units_manifestation_idx ON units (is_manifestation);
