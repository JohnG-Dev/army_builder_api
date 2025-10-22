-- name: GetAbilitiesForUnit :many
SELECT *
FROM abilities
WHERE unit_id = $1
ORDER BY phase ASC, name ASC;

-- name: GetAbilitiesForFaction :many
SELECT *
FROM abilities
WHERE faction_id = $1
ORDER BY phase ASC, name ASC;

-- name: GetAbilityByID :one
SELECT *
FROM abilities
WHERE id = $1;

-- name: GetAbilitiesByType :many
SELECT *
FROM abilities
WHERE type = $1
ORDER BY name ASC;

-- name: GetAbilitiesByPhase :many
SELECT *
FROM abilities
WHERE phase = $1
ORDER BY name ASC;

-- name: GetAllAbilities :many
SELECT *
FROM abilities
ORDER BY unit_id, faction_id, phase ASC, name ASC;

-- name: CreateAbility :one
INSERT INTO abilities (unit_id, faction_id, name, description, type, phase, version, source)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: UpdateAbility :one
UPDATE abilities
SET name = $2, description = $3, type = $4, phase = $5, version = $6, source = $7, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteAbility :exec
DELETE FROM abilities
WHERE id = $1;
