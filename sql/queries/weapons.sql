-- name: GetWeaponsForUnit :many
SELECT * FROM weapons 
WHERE unit_id = $1
ORDER BY unit_id, name ASC;

-- name: GetWeaponByID :one
SELECT * FROM weapons 
WHERE id = $1;

-- name: CreateWeapon :one
INSERT INTO weapons (unit_id, name, range, attacks, to_hit, to_wound, rend, damage)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: DeleteWeapon :exec
DELETE FROM weapons 
WHERE id = $1;
