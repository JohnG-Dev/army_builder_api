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
