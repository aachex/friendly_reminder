package sqlite

import (
	"crypto/rand"
	"database/sql"
	"encoding/base32"
	"sync"

	"github.com/artemwebber1/friendly_reminder/internal/models"
)

type UnverifiedUsersRepository struct {
	mu sync.Mutex
	db *sql.DB
}

func NewUnverifiedUsersRepository(db *sql.DB) *UnverifiedUsersRepository {
	return &UnverifiedUsersRepository{
		db: db,
		mu: sync.Mutex{},
	}
}

// TokenExists возвращает true если существует указанный токен для подтверждения электронной почты.
func (r *UnverifiedUsersRepository) TokenExists(t string) bool {
	row := r.db.QueryRow("SELECT token FROM unverified_users WHERE token = $1", t)
	return row.Scan() != sql.ErrNoRows
}

// CreateToken добавляет пользователя в базу данных, как не подтвердившего электронную почту, и создаёт токен для подтверждения.
// Возвращает сам токен и ошибку.
func (r *UnverifiedUsersRepository) CreateToken(email, pwd string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	token := generateToken()

	_, err := r.db.Exec("INSERT INTO unverified_users(user_email, user_password, token) VALUES($1, $2, $3)", email, pwd, token)
	if err != nil {
		return "", err
	}

	return token, err
}

// DeleteToken удаляет токен из базы данных.
func (r *UnverifiedUsersRepository) DeleteToken(token string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, err := r.db.Exec("DELETE FROM unverified_users WHERE token = $1", token)
	return err
}

// UpdateToken создаёт новый токен для пользователя с указанным email.
func (r *UnverifiedUsersRepository) UpdateToken(email string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	token := generateToken()

	_, err := r.db.Exec("UPDATE unverified_users SET token = $1 WHERE user_email = $2", token, email)
	if err != nil {
		return "", err
	}

	return token, nil
}

// HasToken возвращает true, если для указанной электронной почты уже сгенерирован токен.
func (r *UnverifiedUsersRepository) HasToken(email string) bool {
	row := r.db.QueryRow("SELECT token FROM unverified_users WHERE user_email = $1", email)
	return row.Scan() != sql.ErrNoRows
}

// GetUserByToken получает пользователя по токену.
func (r *UnverifiedUsersRepository) GetUserByToken(token string) (models.User, error) {
	var user models.User

	row := r.db.QueryRow("SELECT user_email, user_password FROM unverified_users WHERE token = $1", token)

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
