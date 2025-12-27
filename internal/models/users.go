package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/geekilx/restaurantAPI/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64     `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type UserModel struct {
	DB *sql.DB
}

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

func (p *password) Matches(plainpassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hashPassword, []byte(plainpassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func (m *UserModel) Insert(user *User) error {
	stmt := `INSERT INTO users (first_name, last_name, email, password_hash) VALUES($1, $2, $3, $4) RETURNING id, created_at, is_active`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{user.FirstName, user.LastName, user.Email, user.Password.hashPassword}

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

func ValidateUsers(v *validator.Validator, user User, password string) {
	v.Check(user.FirstName == "", "firstName", "you have to provide first name")
	v.Check(user.LastName == "", "lastName", "you have to provide last name")
	v.Check(len(user.FirstName) > 50, "firstName", "first name must be greater than 0 and lees than 50 characters")
	v.Check(len(user.LastName) > 50, "lastName", "last name must be greater than 0 and lees than 50 characters")
	v.Check(!validator.CheckEmail(user.Email, validator.EmailRX), "email", "you should provide a valid email address")
	v.Check(len(password) < 6, "password", "passwor dmust be greater than 6 characters")
}
