CREATE TABLE units (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  faction_id UUID NOT NULL REFERENCES factions(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  points INT,
  move TEXT,
  health INT,
  save TEXT,
  ward TEXT,
  control INT,
  rend TEXT,
  attacks TEXT,
  damage TEXT,
  summon_cost TEXT,
  banishment TEXT,
  is_manifestation BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);
