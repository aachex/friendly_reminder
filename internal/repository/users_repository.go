package repository

import (
	"context"
	"database/sql"
)

type UsersRepository interface {
	// AddUser добавляет нового пользователя.
	//
	// Возвращает id нового пользователя и ошибку.
	AddUser(ctx context.Context, email, password string) (int64, error)

	// DeleteUser удаляет пользователя из базы данных.
	DeleteUser(ctx context.Context, email string) error

	// Subscribe подписывает пользователя на рассылку электронных писем.
	// Если параметр subscribe = true, пользователь будет подписан на рассылку, иначе будет отписан.
	Subscribe(ctx context.Context, email string, subscr bool) error

	// GetEmailsSubscribed возвращает список зарегестрированных электронных почт пользователей, подписанных на рассылку.
	GetEmailsSubscribed(ctx context.Context) ([]string, error)

	// EmailExists возвращает true если пользователь с данной электронной почтой уже существует.
	EmailExists(ctx context.Context, email string) bool

	// UserExists возвращает true, если существует пользователь с указанной почтой и паролем.
	UserExists(ctx context.Context, email, password string) bool
}

type usersRepository struct {
	db *sql.DB
}

func NewUsersRepository(db *sql.DB) UsersRepository {
	return &usersRepository{
		db: db,
	}
}

// AddUser добавляет нового пользователя.
//
// Возвращает id нового пользователя и ошибку.
func (r *usersRepository) AddUser(ctx context.Context, email, password string) (id int64, err error) {
	res, err := r.db.ExecContext(ctx, "INSERT INTO users(email, password) VALUES($1, $2)", email, password)
	if err != nil {
		return -1, err
	}
	id, err = res.LastInsertId()
	if err != nil {
		return -1, err
	}
	return id, err
}

// DeleteUser удаляет пользователя из базы данных.
func (r *usersRepository) DeleteUser(ctx context.Context, email string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE email = $1", email)
	return err
}

// Subscribe подписывает пользователя на рассылку электронных писем.
// Если параметр subscribe = true, пользователь будет подписан на рассылку, иначе будет отписан.
func (r *usersRepository) Subscribe(ctx context.Context, email string, subscribe bool) error {
	_, err := r.db.ExecContext(ctx, "UPDATE users SET subscribed = $1 WHERE email = $2", subscribe, email)
	return err
}

// GetEmailsSubscribed возвращает список зарегестрированных электронных почт пользователей, подписанных на рассылку.
func (r *usersRepository) GetEmailsSubscribed(ctx context.Context) (emails []string, err error) {
	rows, err := r.db.QueryContext(ctx, "SELECT email FROM users WHERE subscribed = 1")
	if err != nil {
		return nil, err
	}

	var email string

	for rows.Next() {
		err = rows.Scan(&email)
		if err != nil {
			return nil, err
		}
		emails = append(emails, email)
	}

	return emails, nil
}

// EmailExists возвращает true если пользователь с данной электронной почтой уже существует.
func (r *usersRepository) EmailExists(ctx context.Context, email string) bool {
	row := r.db.QueryRowContext(ctx, "SELECT email FROM users WHERE email = $1", email)
	return row.Scan() != sql.ErrNoRows
}

// UserExists возвращает true, если существует пользователь с указанной почтой и паролем.
func (r *usersRepository) UserExists(ctx context.Context, email, password string) bool {
	row := r.db.QueryRowContext(ctx, "SELECT email, password FROM users WHERE email = $1 AND password = $2", email, password)
	return row.Scan() != sql.ErrNoRows
}
