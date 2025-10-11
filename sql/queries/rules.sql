-- name: GetRulesForGame :many
SELECT * FROM rules WHERE game_id = $1 ORDER BY created_at DESC;

-- name: GetRulesByType :many
SELECT * FROM rules WHERE rule_type = $1 ORDER BY created_at DESC;

-- name: GetRuleByID :one
SELECT * FROM rules WHERE id = $1;

-- name: CreateRule :one
INSERT INTO rules (game_id, name, description, rule_type)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: DeleteRule :exec
DELETE FROM rules WHERE id = $1;
