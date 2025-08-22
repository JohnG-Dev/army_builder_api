-- name: GetUnits :many
SELECT * FROM units WHERE faction_id = $1 ORDER BY created_at DESC;

-- name: GetUnit :one
SELECT * FROM units WHERE id = $1;

-- name: CreateUnit :one
INSERT INTO units (faction_id, name, points, move, health, save, ward, control)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: DeleteUnit :exec
DELETE FROM units WHERE id = $1;
