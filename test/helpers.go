package test

import "database/sql"

func cleanDb(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM Users; DELETE FROM Items")
	return err
}
