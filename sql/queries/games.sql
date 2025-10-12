-- name: GetGames :many
SELECT * FROM games
ORDER BY name ASC;

-- name: GetGame :one
SELECT * FROM games 
WHERE id = $1;

-- name: CreateGame :one
INSERT INTO games (name, edition)
VALUES ($1, $2)
RETURNING *;

-- name: DeleteGame :exec
DELETE FROM games 
WHERE id = $1;
