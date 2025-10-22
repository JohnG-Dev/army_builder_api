CREATE TABLE units (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  faction_id UUID NOT NULL REFERENCES factions(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  is_manifestation BOOLEAN NOT NULL DEFAULT false,
  
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
  min_size INT NOT NULL DEFAULT 1,
  max_size INT NOT NULL DEFAULT 999,
  
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
