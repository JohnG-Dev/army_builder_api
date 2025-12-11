-- name: GetEnhancements :many
SELECT *
FROM enhancements
ORDER BY faction_id, name ASC;

-- name: GetEnhancementsForFaction :many
SELECT *
FROM enhancements
WHERE faction_id = $1
ORDER BY name ASC;

-- name: GetEnhancementByID :one
SELECT *
FROM enhancements
WHERE id = $1;

-- name: GetEnhancementsByType :many
SELECT *
FROM enhancements
WHERE enhancement_type = $1
ORDER BY faction_id, name ASC;

-- name: CreateEnhancement :one
INSERT INTO enhancements (faction_id, name, enhancement_type, description, points, is_unique, version, source)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: UpdateEnhancement :one
UPDATE enhancements
SET name = $2, enhancement_type = $3, description = $4, points = $5, version = $6, source = $7, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteEnhancement :exec
DELETE FROM enhancements
WHERE id = $1;
