package test

import (
	"database/sql"

	"github.com/artemwebber1/friendly_reminder/internal/models"
)

func cleanDb(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM Users; DELETE FROM Items")
	return err
}

func equalSlices(s1, s2 []models.ListItem) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}
