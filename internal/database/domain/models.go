package domain

import (
	"database/sql"
	"errors"
)

var ErrRecordNotFound = errors.New("record not found")

type Models struct {
	App     AppModel
	Users   UserModel
	Courses CourseModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		App:     AppModel{DB: db},
		Users:   UserModel{DB: db},
		Courses: CourseModel{DB: db},
	}
}
