-- name: CreateUser :one
INSERT INTO users(id, created_at, updated_at, name)
VALUES (
  $1,
  $2,
  $3,
  $4
)
RETURNING *;

-- name: GetUserByName :one
SELECT * FROM users WHERE name LIKE $1;

-- name: GetUserByUuid :one
SELECT * FROM users WHERE id = $1;

-- name: DeleteAll :exec
DELETE FROM users;

-- name: AllUsers :many
SELECT * FROM users;
