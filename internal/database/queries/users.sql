-- name: GetUserById :one
select * from users
where id = $1 limit 1;

-- name: GetUserByEmail :one
select * from users
where email = $1 limit 1;

-- name: CreateUser :one
INSERT INTO users (
  email, password_hash, provider, provider_id
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;