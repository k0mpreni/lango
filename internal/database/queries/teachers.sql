-- name: GetTeacherById :one
select * from teachers
where id = $1 limit 1;

-- name: GetTeacherByUserId :one
select * from teachers
where user_id = $1 limit 1;

-- name: CreateTeacher :one
INSERT INTO teachers (
  user_id, description, picture
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: DeleteTeacher :exec
DELETE FROM teachers
WHERE id = $1;