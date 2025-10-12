-- name: GetAbilitiesForUnit :many
SELECT * FROM abilities 
WHERE unit_id = $1 
ORDER BY phase ASC, name ASC;

-- name: GetAbilityByID :one
SELECT * FROM abilities 
WHERE id = $1;

-- name: GetAbilitiesByType :many
SELECT * FROM abilities WHERE type = $1 
ORDER BY name ASC;

-- name: GetAbilitiesByPhase :many
SELECT * FROM abilities 
WHERE phase = $1 
ORDER BY name ASC;

-- name: CreateAbility :one
INSERT INTO abilities (unit_id, faction_id, name, description, type, phase)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: DeleteAbility :exec
DELETE FROM abilities 
WHERE id = $1;
