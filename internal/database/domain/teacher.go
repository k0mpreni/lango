package domain

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Teacher struct {
	Id          uuid.UUID
	UserId      uuid.UUID
	Description string
	Picture     string
	// Topics
}

type TeacherModel struct {
	DB *sql.DB
}

func (m *TeacherModel) GetById(id uuid.UUID) (*Teacher, error) {
	// err := supa.Client.DB.From("teachers").
	// 	Select("*").
	// 	Single().
	// 	Eq("user_id", id.String()).
	// 	Execute(&teacher)

	query := `SELECT id, user_id, description, picture FROM teachers WHERE user_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var teacher Teacher

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&teacher.Id,
		&teacher.UserId,
		&teacher.Description,
		&teacher.Picture,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &teacher, nil
}

func (m *TeacherModel) Create(teacher *Teacher) error {
	query := `
        insert into teachers (user_id, description, picture)
        values ($1, $2, $3)
        returning id 
  `

	args := []any{teacher.UserId, teacher.Description, teacher.Picture}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).
		Scan(&teacher.Id)
	if err != nil {
		return err
	}
	return nil
}
