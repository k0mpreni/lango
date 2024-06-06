package domain

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Student struct {
	Id     uuid.UUID
	UserId uuid.UUID
}

type StudentModel struct {
	DB *sql.DB
}

func (m *StudentModel) GetById(id uuid.UUID) (*Student, error) {
	// err := supa.Client.DB.From("students").
	// 	Select("*").
	// 	Single().
	// 	Eq("user_id", id.String()).
	// 	Execute(&student)

	query := `SELECT id FROM students WHERE user_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var student Student

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&student.Id,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &student, nil
}

func (m *StudentModel) Create(student *Student) error {
	query := `
        insert into students (user_id)
        values ($1)
        returning id 
  `

	args := []any{student.UserId}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).
		Scan(&student.Id)
	if err != nil {
		return err
	}
	return nil
}
