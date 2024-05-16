package domain

import "database/sql"

type Models struct {
	Users UserModel
	App   AppModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users: UserModel{DB: db},
		App:   AppModel{DB: db},
	}
}
