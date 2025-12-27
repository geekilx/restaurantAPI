package models

import (
	"database/sql"
	"errors"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email address")
)

type Models struct {
	UserModel *UserModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		UserModel: &UserModel{DB: db},
	}
}
