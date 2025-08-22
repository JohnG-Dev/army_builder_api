-- name: GetRulesForGame :many
SELECT * FROM rules WHERE game_id = $1 ORDER BY created_at DESC;

-- name: GetRule :one
SELECT * FROM rules WHERE id = $1;

-- name: CreateRule :one
INSERT INTO rules (game_id, name, description)
VALUES ($1, $2, $3)
RETURNING *;

-- name: DeleteRule :exec
DELETE FROM rules WHERE id = $1;
