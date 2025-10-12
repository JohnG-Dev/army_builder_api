-- name: GetUnits :many
SELECT * FROM units 
WHERE faction_id = $1 
ORDER BY name ASC;

-- name: GetUnitByID :one
SELECT * FROM units 
WHERE id = $1;

-- name: ListUnits :many
SELECt * FROM units;

-- name: GetNonManifestationUnits :many
SELECT * FROM units
WHERE is_manifestation = FALSE
ORDER BY name ASC;

-- name: GetManifestationByID :one
SELECT * FROM units
WHERE id = $1 AND is_manifestation = TRUE;

-- name: GetManifestations :many
SELECT * FROM units
WHERE is_manifestation = TRUE
ORDER BY name ASC;

-- name: CreateUnit :one
INSERT INTO units (faction_id, name, points, move, health, save, ward, control, min_size, max_size)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: DeleteUnit :exec
DELETE FROM units 
WHERE id = $1;
