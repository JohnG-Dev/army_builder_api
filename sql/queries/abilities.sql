-- name: GetAbilitiesForUnit :many
SELECT * FROM abilities WHERE unit_id = $1 ORDER BY created_at DESC;

-- name: GetAbility :one
SELECT * FROM abilities WHERE id = $1;

-- name: CreateAbility :one
INSERT INTO abilities (unit_id, name, description)
VALUES ($1, $2, $3)
RETURNING *;

-- name: DeleteAbility :exec
DELETE FROM abilities WHERE id = $1;
