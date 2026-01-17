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
INSERT INTO factions (game_id, name, allegiance, version, source, is_army_of_renown, is_regiment_of_renown, parent_faction_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: UpdateFaction :one
UPDATE factions
SET name = $2, allegiance = $3, version = $4, source = $5, is_army_of_renown = $6, is_regiment_of_renown = $7, parent_faction_id = $8, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteFaction :exec
DELETE FROM factions
WHERE id = $1;
