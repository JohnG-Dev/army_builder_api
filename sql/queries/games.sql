-- name: GetGames :many
SELECT *
FROM games
ORDER BY name ASC;

-- name: GetGame :one
SELECT *
FROM games
WHERE id = $1;

-- name: GetGameByName :one
SELECT *
FROM games
WHERE name = $1;

-- name: CreateGame :one
INSERT INTO games (name, edition, version, source)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateGame :one
UPDATE games
SET edition = $2, version = $3, source = $4, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteGame :exec
DELETE FROM games
WHERE id = $1;
