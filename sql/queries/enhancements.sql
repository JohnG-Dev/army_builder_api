-- name: GetEnhancementsForFaction :many
SELECT * FROM enhancements
WHERE faction_id = $1
ORDER BY faction_id, name ASC;

-- name: GetEnhancementByID :one
SELECT * FROM enhancements 
WHERE id = $1;

-- name: CreateEnhancement :one
INSERT INTO enhancements (faction_id, name, description, points)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: DeleteEnhancement :exec
DELETE FROM enhancements 
WHERE id = $1;



