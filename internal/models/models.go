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
	ErrRestaurantNotFound      = errors.New("no restaurant found")
)

type Models struct {
	Users       *UserModel
	Restaurants *RestaurantModel
	Tokens      *TokenModel
	Permissions *PermissionModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users:       &UserModel{DB: db},
		Restaurants: &RestaurantModel{DB: db},
		Tokens:      &TokenModel{DB: db},
		Permissions: &PermissionModel{DB: db},
	}
}
