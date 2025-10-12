CREATE TABLE abilities (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  unit_id UUID REFERENCES units(id) ON DELETE CASCADE,
  faction_id UUID REFERENCES factions(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  description TEXT,
  type TEXT,    -- 'passive', 'spell', 'prayer', 'trait'
  phase TEXT,   -- 'hero', 'movement', 'charge', 'combat', 'end_of_turn'
  version TEXT,
  source TEXT,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);
