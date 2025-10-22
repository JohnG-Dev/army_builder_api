-- name: GetAbilityEffectsForAbility :many
SELECT *
FROM ability_effects
WHERE ability_id = $1
ORDER BY stat ASC;

-- name: GetAbilityEffectByID :one
SELECT *
FROM ability_effects
WHERE id = $1;

-- name: GetAllAbilityEffects :many
SELECT *
FROM ability_effects
ORDER BY ability_id, stat ASC;

-- name: CreateAbilityEffect :one
INSERT INTO ability_effects (ability_id, stat, modifier, condition, description, version, source)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateAbilityEffect :one
UPDATE ability_effects
SET stat = $2, modifier = $3, condition = $4, description = $5, version = $6, source = $7, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteAbilityEffect :exec
DELETE FROM ability_effects
WHERE id = $1;
