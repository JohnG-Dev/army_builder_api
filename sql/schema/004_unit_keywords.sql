CREATE TABLE unit_keywords (
  unit_id UUID NOT NULL REFERENCES units(id) ON DELETE CASCADE,
  keyword_id UUID NOT NULL REFERENCES keywords(id) ON DELETE CASCADE,
  value TEXT NOT NULL DEFAULT '',
  PRIMARY KEY (unit_id, keyword_id)
);
