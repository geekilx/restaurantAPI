package models

import (
	"database/sql"
	"errors"
)

var (
	ErrDuplicateEmail          = errors.New("duplicate email address")
	ErrRecordNotFound          = errors.New("no record found")
	ErrConflictEdit            = errors.New("conflict edit")
	ErrDuplicateRestaurantName = errors.New("duplicate restaurant name")
)

type Models struct {
	Users       *UserModel
	Restaurants *RestaurantModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users:       &UserModel{DB: db},
		Restaurants: &RestaurantModel{DB: db},
	}
}
