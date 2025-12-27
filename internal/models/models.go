package models

import (
	"database/sql"
	"errors"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email address")
	ErrRecordNotFound = errors.New("no record found")
	ErrConflictEdit   = errors.New("conflict edit")
)

type Models struct {
	Users *UserModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users: &UserModel{DB: db},
	}
}
