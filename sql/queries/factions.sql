-- name: GetFactions :many
SELECT * FROM factions WHERE game_id = $1 ORDER BY created_at DESC;

-- name: CreateFaction :one
INSERT INTO factions (game_id, name, allegiance)
VALUES ($1, $2, $3)
RETURNING *;
