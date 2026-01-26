package models

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/geekilx/restaurantAPI/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int64     `json:"id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Email        string    `json:"email"`
	Role         string    `json:"role"`
	RestaurantID *int64    `json:"restaurant_id"`
	Password     password  `json:"-"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
}

type UserModel struct {
	DB *sql.DB
}

var AnonymousUser = &User{}

type password struct {
	plainPassword *string
	hashPassword  []byte
}

func (p *password) Set(plainpassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainpassword), 12)
	if err != nil {
		return err
	}

	p.hashPassword = hash
	p.plainPassword = &plainpassword
	return nil
}

func (p *password) Matches(plainPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hashPassword, []byte(plainPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, err
		default:
			return false, err
		}
	}

	return true, nil
}

func IsAnonymous(user *User) bool {
	return user == AnonymousUser
}

func (m *UserModel) Insert(user *User) error {
	stmt := `INSERT INTO users (first_name, last_name, email, password_hash, role) VALUES($1, $2, $3, $4, $5) RETURNING id, created_at, is_active`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{user.FirstName, user.LastName, user.Email, user.Password.hashPassword, user.Role}

	err := m.DB.QueryRowContext(ctx, stmt, args...).Scan(&user.ID, &user.CreatedAt, &user.IsActive)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil

}

func (m *UserModel) GetUser(id int64) (*User, error) {

	if id < 1 {
		return nil, ErrRecordNotFound
	}

	stmt := `SELECT id, first_name, last_name, email, created_at, is_active, role, restaurant_id FROM users WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User

	err := m.DB.QueryRowContext(ctx, stmt, id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.CreatedAt, &user.IsActive, &user.Role, &user.RestaurantID)
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

func (m *UserModel) Update(user *User) error {
	stmt := `UPDATE users SET first_name = $1, last_name = $2, email = $3, password_hash = $4, is_active = $5, restaurant_id = $6, last_updated = NOW() where id = $7`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{user.FirstName, user.LastName, user.Email, user.Password.hashPassword, user.IsActive, user.RestaurantID, user.ID}

	result, err := m.DB.ExecContext(ctx, stmt, args...)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "users_email_key"):
			return ErrDuplicateEmail
		}
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (m *UserModel) Delete(id int64) error {

	stmt := `DELETE FROM users WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	rowsAffected, err := rows.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (m *UserModel) GetUserByEmail(email string) (*User, error) {
	stmt := `SELECT id, password_hash FROM users WHERE email = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User

	err := m.DB.QueryRowContext(ctx, stmt, email).Scan(&user.ID, &user.Password.hashPassword)
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

func ValidateUsers(v *validator.Validator, user User, password string) {
	v.Check(user.FirstName == "", "firstName", "you have to provide first name")
	v.Check(user.LastName == "", "lastName", "you have to provide last name")
	v.Check(len(user.FirstName) > 50, "firstName", "first name must be greater than 0 and lees than 50 characters")
	v.Check(len(user.LastName) > 50, "lastName", "last name must be greater than 0 and lees than 50 characters")
	v.Check(!validator.CheckEmail(user.Email, validator.EmailRX), "email", "you should provide a valid email address")
	v.Check(len(password) < 6, "password", "password must be greater than 6 characters")
}
