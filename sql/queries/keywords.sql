-- name: GetAllKeywords :many
SELECT *
FROM keywords
ORDER BY game_id, name ASC;

-- name: GetKeywordsByGame :many
SELECT *
FROM keywords
WHERE game_id = $1
ORDER BY name ASC;

-- name: GetKeywordByID :one
SELECT *
FROM keywords
WHERE id = $1;

-- name: GetKeywordsForUnit :many
SELECT k.*
FROM keywords k
JOIN unit_keywords uk ON k.id = uk.keyword_id
WHERE uk.unit_id = $1
ORDER BY k.name ASC;

-- name: GetUnitsWithKeyword :many
SELECT DISTINCT u.*
FROM units u
JOIN unit_keywords uk ON u.id = uk.unit_id
JOIN keywords k ON uk.keyword_id = k.id
WHERE k.name = $1
ORDER BY u.name ASC;

-- name: GetUnitsWithKeywordAndValue :many
SELECT DISTINCT u.*
FROM units u
JOIN unit_keywords uk ON u.id = uk.unit_id
JOIN keywords k ON uk.keyword_id = k.id
WHERE k.name = $1 AND uk.value = $2
ORDER BY u.name ASC;

-- name: CreateKeyword :one
INSERT INTO keywords (game_id, name, description, version, source)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: AddKeywordToUnit :exec
INSERT INTO unit_keywords (unit_id, keyword_id, value)
VALUES ($1, $2, $3)
ON CONFLICT DO NOTHING;

-- name: RemoveKeywordFromUnit :exec
DELETE FROM unit_keywords
WHERE unit_id = $1 AND keyword_id = $2;

-- name: UpdateKeywordValue :exec
UPDATE unit_keywords
SET value = $3
WHERE unit_id = $1 AND keyword_id = $2;

-- name: DeleteKeyword :exec
DELETE FROM keywords
WHERE id = $1;
