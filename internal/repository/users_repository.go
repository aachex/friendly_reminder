package repository

import (
	"database/sql"
)

type UsersRepository interface {
	// AddUser добавляет нового пользователя.
	//
	// Возвращает id нового пользователя и ошибку.
	AddUser(email, passwordHash string) (int64, error)

	// GetEmails возвращает список зарегестрированных электронных почт.
	GetEmails() ([]string, error)

	// EmailExists возвращает true если пользователь с данной электронной почтой уже существует.
	EmailExists(email string) bool
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
func (r *usersRepository) AddUser(email, passwordHash string) (id int64, err error) {
	res, err := r.db.Exec("INSERT INTO Users(email, password) VALUES($1, $2);", email, passwordHash)
	if err != nil {
		return -1, err
	}
	id, err = res.LastInsertId()
	if err != nil {
		return -1, err
	}
	return id, err
}

// GetEmails возвращает список зарегестрированных электронных почт.
func (r *usersRepository) GetEmails() (emails []string, err error) {
	rows, err := r.db.Query("SELECT email FROM Users;")
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
func (r *usersRepository) EmailExists(email string) bool {
	_, err := r.db.Exec("SELECT email FROM Users WHERE email = $1", email)
	return err != sql.ErrNoRows
}
