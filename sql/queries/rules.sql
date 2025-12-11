-- name: GetAllRules :many
SELECT *
FROM rules
ORDER BY game_id, rule_type ASC, name ASC;

-- name: GetRulesForGame :many
SELECT *
FROM rules
WHERE game_id = $1
ORDER BY rule_type ASC, name ASC;

-- name: GetRulesByType :many
SELECT *
FROM rules
WHERE game_id = $1 AND rule_type = $2
ORDER BY name ASC;

-- name: GetRuleByID :one
SELECT *
FROM rules
WHERE id = $1;

-- name: CreateRule :one
INSERT INTO rules (game_id, name, description, rule_type, text, version, source)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateRule :one
UPDATE rules
SET name = $2, description = $3, rule_type = $4, version = $5, source = $6, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteRule :exec
DELETE FROM rules
WHERE id = $1;
