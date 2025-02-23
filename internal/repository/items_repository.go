package repository

import (
	"database/sql"

	"github.com/artemwebber1/friendly_reminder/internal/models"
)

type ItemsRepository interface {
	// AddItem добавляет новую задачу в список пользователя.
	AddItem(value, userEmail string) error

	// GetList возвращает список дел пользователя с указанным email.
	GetList(userEmail string) ([]models.ListItem, error)
}

type itemsRepository struct {
	db *sql.DB
}

func NewItemsRepository(db *sql.DB) ItemsRepository {
	return &itemsRepository{
		db: db,
	}
}

func (r *itemsRepository) AddItem(value, userEmail string) error {
	_, err := r.db.Exec(`
		INSERT INTO Items(value, numberInList, userEmail) VALUES(
			$1,
			(SELECT COUNT(*) FROM Items WHERE userEmail = $2) + 1,
			$2)`,
		value, userEmail)
	return err
}

// GetList возвращает список дел пользователя с указанным id.
func (r *itemsRepository) GetList(userEmail string) ([]models.ListItem, error) {
	rows, err := r.db.Query("SELECT value, numberInList FROM Items WHERE userEmail = $1 ORDER BY numberInList", userEmail)
	if err != nil {
		return nil, err
	}

	var items []models.ListItem
	var item models.ListItem
	for rows.Next() {
		rows.Scan(&item.Value, &item.NumberInList)
		items = append(items, item)
	}

	return items, nil
}
