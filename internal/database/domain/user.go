package domain

import (
	"database/sql"

	"github.com/google/uuid"
)

type User struct {
	Id    uuid.UUID
	Role  int8
	Email string
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) GetByEmail(email string) (*User, error) {
	return nil, nil
}
