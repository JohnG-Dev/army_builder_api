CREATE TABLE weapons (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  unit_id UUID NOT NULL REFERENCES units(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  range TEXT NOT NULL DEFAULT '',
  attacks TEXT NOT NULL DEFAULT '',
  hit_stats TEXT NOT NULL DEFAULT '',
  wound_strength TEXT NOT NULL DEFAULT '',
  rend_ap TEXT NOT NULL DEFAULT '',
  damage TEXT NOT NULL DEFAULT '',
  version TEXT NOT NULL DEFAULT '',
  source TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX weapons_unit_idx ON weapons (unit_id, name ASC);
