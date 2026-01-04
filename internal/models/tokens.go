package models

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"errors"
	"time"

	"github.com/geekilx/restaurantAPI/internal/validator"
)

var (
	ActivationScope     = "activation"
	AuthenticationScope = "authentication"
)

type Token struct {
	PlainHash string    `json:"token"`
	Hash      []byte    `json:"-"`
	UserID    int64     `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

type TokenModel struct {
	DB *sql.DB
}

func GenerateToken(ttl time.Duration, userID int64, scope string) (*Token, error) {

	token := Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	randomeBytes := make([]byte, 16)

	_, err := rand.Read(randomeBytes)
	if err != nil {
		return &token, err
	}

	token.PlainHash = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomeBytes)

	hash := sha256.Sum256([]byte(token.PlainHash))

	token.Hash = hash[:]

	return &token, nil
}

func (m TokenModel) New(ttl time.Duration, userID int64, scope string) (*Token, error) {
	token, err := GenerateToken(ttl, userID, scope)
	if err != nil {
		return nil, err
	}

	err = m.Insert(token)
	if err != nil {
		return nil, err
	}

	return token, err

}

func (m TokenModel) Insert(token *Token) error {
	stmt := `INSERT INTO tokens (hash, user_id, expiry, scope) VALUES($1, $2, $3, $4)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{token.Hash, token.UserID, token.Expiry, token.Scope}

	_, err := m.DB.ExecContext(ctx, stmt, args...)
	return err

}

func (m *TokenModel) GetByToken(plainToken string) (int64, error) {

	hash := sha256.Sum256([]byte(plainToken))

	stmt := `SELECT * FROM tokens WHERE hash = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int64
	err := m.DB.QueryRowContext(ctx, stmt, hash[:]).Scan(&id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return 0, ErrRecordNotFound
		default:
			return 0, err
		}
	}

	return id, nil

}

func (m *TokenModel) DeleteAllTokenForUser(userID int64, scope string) error {
	stmt := `DELETE FROM tokens WHERE user_id = $1 AND scope = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, stmt, userID, scope)
	return err

}
func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "must be provided")
	v.Check(len(tokenPlaintext) == 26, "token", "must be 26 bytes long")
}
