package repository

import (
	"crypto/rand"
	"database/sql"
	"encoding/base32"

	"github.com/artemwebber1/friendly_reminder/internal/models"
)

type EmailTokensRepository interface {
	// TokenExists возвращает true если существует указанный токен для подтверждения электронной почты.
	TokenExists(t string) bool

	// CreateToken создаёт новый токен для подтверждения электронной почты. Возвращает сам токен и ошибку.
	CreateToken(email, pwd string) (string, error)

	// DeleteToken удаляет токен из базы данных.
	DeleteToken(token string) error

	// UpdateToken создаёт новый токен для пользователя с указанным email.
	UpdateToken(email string) (string, error)

	// HasToken возвращает true, если для указанной электронной почты уже сгенерирован токен.
	HasToken(email string) bool

	GetUserByToken(token string) (models.User, error)
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
func (r *emailTokensRepository) CreateToken(email, pwd string) (string, error) {
	token := generateToken()

	_, err := r.db.Exec("INSERT INTO EmailConfirmTokens(user_email, user_password, token) VALUES($1, $2, $3)", email, pwd, token)
	if err != nil {
		return "", err
	}

	return token, err
}

// DeleteToken удаляет токен из базы данных.
func (r *emailTokensRepository) DeleteToken(token string) error {
	_, err := r.db.Exec("DELETE FROM EmailConfirmTokens WHERE token = $1", token)
	return err
}

func (r *emailTokensRepository) UpdateToken(email string) (string, error) {
	token := generateToken()

	_, err := r.db.Exec("UPDATE EmailConfirmTokens SET token = $1 WHERE user_email = $2", token, email)
	if err != nil {
		return "", err
	}

	return token, nil
}

// HasToken возвращает true, если для указанной электронной почты уже сгенерирован токен.
func (r *emailTokensRepository) HasToken(email string) bool {
	row := r.db.QueryRow("SELECT token FROM EmailConfirmTokens WHERE user_email = $1", email)
	return row.Scan() != sql.ErrNoRows
}

func (r *emailTokensRepository) GetUserByToken(token string) (models.User, error) {
	var user models.User

	row := r.db.QueryRow("SELECT user_email, user_password FROM EmailConfirmTokens WHERE token = $1", token)

	err := row.Scan(&user.Email, &user.Password)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func generateToken() string {
	tokBytes := make([]byte, 32)
	rand.Read(tokBytes)
	return base32.StdEncoding.EncodeToString(tokBytes)
}
