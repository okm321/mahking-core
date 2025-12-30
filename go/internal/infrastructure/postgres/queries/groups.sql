-- name: ListGroups :many
SELECT
  id,
  uid,
  name,
  created_at,
  updated_at
FROM
  groups
ORDER BY id;

-- name: CreateGroup :one
INSERT INTO groups (name)
VALUES (@name)
RETURNING id, uid, name, created_at, updated_at;
