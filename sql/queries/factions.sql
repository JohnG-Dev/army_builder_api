-- name: GetFactionsByID :many
SELECT * FROM factions 
WHERE game_id = $1 
ORDER BY name ASC;

-- name: GetAllFactions :many
SELECT * FROM factions 
ORDER BY game_id, name ASC;

-- name: GetFaction :one
SELECT * FROM factions 
WHERE id = $1;

-- name: CreateFaction :one
INSERT INTO factions (game_id, name, allegiance)
VALUES ($1, $2, $3)
RETURNING *;

-- name: DelteFaction :exec
DELETE FROM factions 
WHERE id = $1;
