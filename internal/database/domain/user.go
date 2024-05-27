package domain

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id           uuid.UUID
	CreatedAt    time.Time
	Role         int8
	Email        string
	UserId       uuid.UUID
	ContactInfos sql.NullString
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) GetById(id uuid.UUID) (*User, error) {
	// err := supa.Client.DB.From("users").
	// 	Select("*").
	// 	Single().
	// 	Eq("user_id", id.String()).
	// 	Execute(&user)

	query := `SELECT id, email, role, created_at, contact_infos FROM users WHERE user_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&user.Id,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
		&user.ContactInfos,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (m *UserModel) GetByEmail(email string) (*User, error) {
	// err = supa.Client.DB.From("users").
	// 	Select("*").
	// 	Single().
	// 	Eq("email", email).
	// 	Execute(&user)

	query := `SELECT id, email, role, created_at, contact_infos FROM users WHERE email = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User

	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&user.Id,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
		&user.ContactInfos,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}
