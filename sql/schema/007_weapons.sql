CREATE TABLE weapons(
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  unit_id UUID NOT NULL REFERENCES units(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  range TEXT,
  attacks TEXT,
  to_hit TEXT,
  to_wound TEXT,
  rend TEXT,
  damage TEXT,
  version TEXT,
  source TEXT,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);
