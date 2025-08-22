-- name: GetFactions :many
SELECT * FROM factions WHERE game_id = $1 ORDER BY created_at DESC;

-- name: GetFaction :one
SELECT * FROM factions where id = $1;

-- name: CreateFaction :one
INSERT INTO factions (game_id, name, allegiance)
VALUES ($1, $2, $3)
RETURNING *;

-- name: DelteFaction :exec
DELETE FROM factions where id = $1;

-- name: GetAllFactions :many
SELECT * FROM factions ORDER BY name;
