CREATE TABLE units (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  faction_id UUID NOT NULL REFERENCES factions(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  points INT NOT NULL,
  move TEXT,
  health INT,
  save TEXT,
  ward TEXT,
  control INT,
  min_size INT DEFAULT 1 NOT NULL,
  max_size INT DEFAULT 1 NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);
