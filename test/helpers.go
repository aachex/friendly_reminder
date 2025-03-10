package test

import (
	"database/sql"
)

const email = "abcde@gmail.com"
const passwordHash = "hashedPassword"

func cleanDb(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM users; DELETE FROM tasks; DELETE FROM unverified_users;")
	return err
}
