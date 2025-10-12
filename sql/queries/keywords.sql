-- name: GetKeywordsForUnit :many
SELECT k.id, k.name, uk.value
FROM keywords k
JOIN unit_keywords uk ON uk.keyword_id = k.id
WHERE uk.unit_id = $1
ORDER BY name ASC;

-- name: GetUnitsWithKeyword :many
SELECT u.*
FROM units u
JOIN unit_keywords uk ON uk.unit_id = u.id
JOIN keywords k ON k.id = uk.keyword_id
WHERE k.name = $1
ORDER BY name ASC;

-- name: GetUnitsWithKeywordAndValue :many
SELECT u.*
FROM units u
JOIN unit_keywords uk ON uk.unit_id = u.id
JOIN keywords k ON k.id = uk.keyword_id
WHERE k.name = $1
  AND uk.value = $2
ORDER BY name ASC;

-- name: AddKeywordToUnit :exec
INSERT INTO unit_keywords (unit_id, keyword_id, value)
VALUES ($1, $2, $3);

-- name: RemoveKeywordFromUnit :exec
DELETE FROM unit_keywords
WHERE unit_id = $1 AND keyword_id = $2;

-- name: CreateKeyword :one
INSERT INTO keywords (name)
VALUES ($1)
RETURNING *;

-- name: GetAllKeywords :many
SELECT * FROM keywords 
ORDER BY name ASC;
