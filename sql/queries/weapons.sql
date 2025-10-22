-- name: GetWeaponsForUnit :many
SELECT *
FROM weapons
WHERE unit_id = $1
ORDER BY name ASC;

-- name: GetWeaponByID :one
SELECT *
FROM weapons
WHERE id = $1;

-- name: GetAllWeapons :many
SELECT *
FROM weapons
ORDER BY unit_id, name ASC;

-- name: CreateWeapon :one
INSERT INTO weapons (unit_id, name, range, attacks, to_hit, to_wound, rend, damage, version, source)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: UpdateWeapon :one
UPDATE weapons
SET name = $2, range = $3, attacks = $4, to_hit = $5, to_wound = $6, rend = $7, damage = $8, version = $9, source = $10, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteWeapon :exec
DELETE FROM weapons
WHERE id = $1;
