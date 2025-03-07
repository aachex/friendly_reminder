package repository

import (
	"crypto/rand"
	"database/sql"
	"encoding/base32"
)

type EmailTokensRepository interface {
	// TokenExists возвращает true если существует указанный токен для подтверждения электронной почты.
	TokenExists(t string) bool

	// CreateToken создаёт новый токен для подтверждения электронной почты. Возвращает сам токен и ошибку.
	CreateToken(email string) (string, error)
}

type emailTokensRepository struct {
	db *sql.DB
}

func NewEmailTokensRepository(db *sql.DB) EmailTokensRepository {
	return &emailTokensRepository{
		db: db,
	}
}

// TokenExists возвращает true если существует указанный токен для подтверждения электронной почты.
func (r *emailTokensRepository) TokenExists(t string) bool {
	row := r.db.QueryRow("SELECT token FROM EmailConfirmTokens WHERE token = $1", t)
	return row.Scan() != sql.ErrNoRows
}

// CreateToken создаёт новый токен для подтверждения электронной почты. Возвращает сам токен и ошибку.
func (r *emailTokensRepository) CreateToken(email string) (string, error) {
	tokBytes := make([]byte, 32)
	rand.Read(tokBytes)
	token := base32.StdEncoding.EncodeToString(tokBytes)

	_, err := r.db.Exec("INSERT INTO EmailConfirmTokens(user_email, token) VALUES($1, $2)", email, token)
	if err != nil {
		return "", err
	}

	return token, err
}
