-- name: GetAllUnits :many
SELECT *
FROM units
WHERE is_manifestation = false
ORDER BY faction_id, name ASC;

-- name: GetUnitsByFaction :many
SELECT *
FROM units
WHERE faction_id = $1 AND is_manifestation = false
ORDER BY name ASC;

-- name: GetUnitByID :one
SELECT *
FROM units
WHERE id = $1;

-- name: GetUnitsByMatchedPlay :many
SELECT *
FROM units
WHERE faction_id = $1 AND matched_play = true AND is_manifestation = false
ORDER BY name ASC;

-- name: GetNonManifestationUnits :many
SELECT *
FROM units
WHERE is_manifestation = false
ORDER BY faction_id, name ASC;

-- name: GetManifestations :many
SELECT *
FROM units
WHERE is_manifestation = true
ORDER BY faction_id, name ASC;

-- name: GetManifestationByID :one
SELECT *
FROM units
WHERE id = $1 AND is_manifestation = true;

-- name: CreateUnit :one
INSERT INTO units (
  faction_id, name, description, is_manifestation,
  move, health, save, ward, control, points,
  summon_cost, banishment,
  min_size, max_size, matched_play, version, source
)
VALUES (
  $1, $2, $3, $4,
  $5, $6, $7, $8, $9, $10,
  $11, $12,
  $13, $14, $15, $16, $17
)
RETURNING *;

-- name: UpdateUnit :one
UPDATE units
SET name = $2, description = $3, move = $4, health = $5, save = $6, ward = $7,
    control = $8, points = $9, summon_cost = $10, banishment = $11,
    min_size = $12, max_size = $13, matched_play = $14, version = $15, source = $16, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteUnit :exec
DELETE FROM units
WHERE id = $1;
