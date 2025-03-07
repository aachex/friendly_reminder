package repository

import "database/sql"

type emailTokensRepository struct {
	db *sql.DB
}

// TokenExists возвращает true если существует указанный токен для подтверждения электронной почты.
func (r *emailTokensRepository) TokenExists(t string) bool {
	row := r.db.QueryRow("SELECT token_val FROM EmailConfirmTokens WHERE token_val = $1", t)
	return row.Scan() != sql.ErrNoRows
}
