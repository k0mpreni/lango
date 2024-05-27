package domain

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
)

type Course struct {
	Id          uuid.UUID
	CreatedAt   time.Time
	TeacherId   uuid.UUID
	StudentId   uuid.UUID
	Title       string
	Description string
	Date        time.Time
	Completed   bool
	Canceled    bool
}

type CourseModel struct {
	DB *sql.DB
}

func (m *CourseModel) Create(course *Course) error {
	query := `
        insert into courses (teacher_id, student_id, title, description, date)
        values ($1, $2, $3, $4, $5)
        returning id, created_at, completed, canceled
  `

	args := []any{course.TeacherId, course.StudentId, course.Title, course.Description, course.Date}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).
		Scan(&course.Id, &course.CreatedAt, &course.Completed, &course.Canceled)
	if err != nil {
		return err
	}
	return nil
}

func (m *CourseModel) Update(id uuid.UUID) (*Course, error) {
	return nil, nil
}

func (m *CourseModel) Delete(id uuid.UUID) (*Course, error) {
	return nil, nil
}

func (m *CourseModel) GetById(id uuid.UUID) (*Course, error) {
	query := `SELECT * FROM courses WHERE id = $1;`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var course Course

	err := m.DB.QueryRowContext(ctx, query, id).Scan()
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &course, nil
}

func (m *CourseModel) GetAllByTeacher(user *User) (*[]Course, error) {
	query := `SELECT * FROM courses WHERE teacher_id = $1;`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var courses []Course

	rows, err := m.DB.QueryContext(ctx, query, user.Id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	for rows.Next() {
		err = rows.Scan(&courses)
		if err != nil {
			log.Fatalf("Failed to retrieve row because %s", err)
		}
	}
	return &courses, nil
}

func (m *CourseModel) GetAllByStudent(user *User) (*[]Course, error) {
	query := `SELECT * FROM courses WHERE student_id = $1;`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var courses []Course
	rows, err := m.DB.QueryContext(ctx, query, user.Id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	for rows.Next() {
		err = rows.Scan(&courses)
		if err != nil {
			log.Fatalf("Failed to retrieve row because %s", err)
		}
	}

	return &courses, nil
}
