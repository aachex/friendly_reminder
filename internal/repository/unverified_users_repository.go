package repository

import (
	"crypto/rand"
	"database/sql"
	"encoding/base32"
	"sync"

	"github.com/artemwebber1/friendly_reminder/internal/models"
)

// UnverifiedUsersRepository является репозиторием неверифицированных пользователей.
//
// Неверифицированный пользователь - это пользователь, который
// регистрировался в системе, но не подтвердил электронную почту.
type UnverifiedUsersRepository interface {
	// TokenExists возвращает true если существует указанный токен для подтверждения электронной почты.
	TokenExists(t string) bool

	// CreateToken добавляет пользователя в базу данных, как не подтвердившего электронную почту, и создаёт токен для подтверждения.
	// Возвращает сам токен и ошибку.
	CreateToken(email, pwd string) (string, error)

	// DeleteToken удаляет токен из базы данных.
	DeleteToken(token string) error

	// UpdateToken создаёт новый токен для пользователя с указанным email.
	UpdateToken(email string) (string, error)

	// HasToken возвращает true, если для указанной электронной почты уже сгенерирован токен.
	HasToken(email string) bool

	// GetUserByToken получает пользователя по токену.
	GetUserByToken(token string) (models.User, error)
}

type unverifiedUsersRepository struct {
	mu sync.Mutex
	db *sql.DB
}

func NewUnverifiedUsersRepository(db *sql.DB) UnverifiedUsersRepository {
	return &unverifiedUsersRepository{
		db: db,
		mu: sync.Mutex{},
	}
}

// TokenExists возвращает true если существует указанный токен для подтверждения электронной почты.
func (r *unverifiedUsersRepository) TokenExists(t string) bool {
	row := r.db.QueryRow("SELECT token FROM unverified_users WHERE token = $1", t)
	return row.Scan() != sql.ErrNoRows
}

// CreateToken добавляет пользователя в базу данных, как не подтвердившего электронную почту, и создаёт токен для подтверждения.
// Возвращает сам токен и ошибку.
func (r *unverifiedUsersRepository) CreateToken(email, pwd string) (string, error) {
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
func (r *unverifiedUsersRepository) DeleteToken(token string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, err := r.db.Exec("DELETE FROM unverified_users WHERE token = $1", token)
	return err
}

// UpdateToken создаёт новый токен для пользователя с указанным email.
func (r *unverifiedUsersRepository) UpdateToken(email string) (string, error) {
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
func (r *unverifiedUsersRepository) HasToken(email string) bool {
	row := r.db.QueryRow("SELECT token FROM unverified_users WHERE user_email = $1", email)
	return row.Scan() != sql.ErrNoRows
}

// GetUserByToken получает пользователя по токену.
func (r *unverifiedUsersRepository) GetUserByToken(token string) (models.User, error) {
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
