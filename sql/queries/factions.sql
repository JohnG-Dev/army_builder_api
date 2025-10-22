-- name: GetAllFactions :many
SELECT *
FROM factions
ORDER BY game_id, name ASC;

-- name: GetFactionsByID :many
SELECT *
FROM factions
WHERE game_id = $1
ORDER BY name ASC;

-- name: GetFaction :one
SELECT *
FROM factions
WHERE id = $1;

-- name: GetFactionsByName :many
SELECT *
FROM factions
WHERE name ILIKE $1
ORDER BY game_id, name ASC;

-- name: CreateFaction :one
INSERT INTO factions (game_id, name, allegiance, version, source)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateFaction :one
UPDATE factions
SET name = $2, allegiance = $3, version = $4, source = $5, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteFaction :exec
DELETE FROM factions
WHERE id = $1;
