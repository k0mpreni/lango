package domain

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID
	CreatedAt time.Time
	Role      int8
	Email     string
	UserId    uuid.UUID
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) GetByEmail(email string) (*User, error) {
	return nil, nil
}
