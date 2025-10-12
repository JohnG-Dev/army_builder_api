-- name: GetBattleFormationsForFaction :many
SELECT * FROM battle_formations
WHERE faction_id = $1
ORDER BY name ASC;

-- name: GetBattleFormationByID :one
SELECT * FROM battle_formations 
WHERE id = $1;

-- name: CreateBattleFormation :one
INSERT INTO battle_formations (game_id, faction_id, name, description, version, source)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: DeleteBattleFormation :exec
DELETE FROM battle_formations 
WHERE id = $1;
