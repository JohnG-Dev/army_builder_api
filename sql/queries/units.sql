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
  faction_id, name, description, is_manifestation, is_unique,
  move, health_wounds, save_stats, ward_fnp, invuln_save, 
  control_oc, toughness, leadership_bravery, points,
  additional_stats,
  summon_cost, banishment,
  min_unit_size, max_unit_size, matched_play, version, source
)
VALUES (
  $1, $2, $3, $4, $5,
  $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
  $16, $17,
  $18, $19, $20, $21, $22
)
RETURNING *;

-- name: UpdateUnit :one
UPDATE units
SET name = $2, description = $3, move = $4, health_wounds = $5, save_stats = $6, 
    ward_fnp = $7, invuln_save = $8, control_oc = $9, toughness = $10, 
    leadership_bravery = $11, points = $12, additional_stats = $13,
    summon_cost = $14, banishment = $15, min_unit_size = $16, max_unit_size = $17, 
    matched_play = $18, version = $19, source = $20, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteUnit :exec
DELETE FROM units
WHERE id = $1;
