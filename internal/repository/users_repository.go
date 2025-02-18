package repository

import (
	"database/sql"
)

type UsersRepository interface {
	// AddUser добавляет нового пользователя.
	AddUser(email, passwordHash string) error

	// GetEmails возвращает список зарегестрированных электронных почт.
	GetEmails() ([]string, error)
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
func (r *usersRepository) AddUser(email, passwordHash string) error {
	_, err := r.db.Exec("INSERT INTO Users(email, passwordHash) VALUES($1, $2);", email, passwordHash)
	return err
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
