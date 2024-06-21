-- name: GetCourse :one
select * from courses
where id = $1 limit 1;

-- name: ListCoursesByTeacherId :many
select * from courses
where teacher_id = $1;

-- name: CreateCourse :one
INSERT INTO courses (
  title, teacher_id, student_id, description, date
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: DeleteCourse :exec
DELETE FROM courses
WHERE id = $1;