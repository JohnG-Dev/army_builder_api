-- name: GetBattleFormationsForGame :many
SELECT *
FROM battle_formations
WHERE game_id = $1
ORDER BY faction_id, name ASC;

-- name: GetBattleFormationsForFaction :many
SELECT *
FROM battle_formations
WHERE faction_id = $1
ORDER BY name ASC;

-- name: GetBattleFormationByID :one
SELECT *
FROM battle_formations
WHERE id = $1;

-- name: GetAllBattleFormations :many
SELECT *
FROM battle_formations
ORDER BY game_id, faction_id, name ASC;

-- name: CreateBattleFormation :one
INSERT INTO battle_formations (game_id, faction_id, name, description, version, source)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateBattleFormation :one
UPDATE battle_formations
SET name = $2, description = $3, version = $4, source = $5, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteBattleFormation :exec
DELETE FROM battle_formations
WHERE id = $1;
