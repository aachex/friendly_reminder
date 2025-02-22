package repository

import (
	"database/sql"

	"github.com/artemwebber1/friendly_reminder/internal/models"
)

type ItemsRepository interface {
	// AddItem добавляет новую задачу в список пользователя.
	AddItem(value string, userId int64) error

	// GetList возвращает список дел пользователя с указанным id.
	GetList(userId int64) ([]models.ListItem, error)
}

type itemsRepository struct {
	db *sql.DB
}

func NewItemsRepository(db *sql.DB) ItemsRepository {
	return &itemsRepository{
		db: db,
	}
}

func (r *itemsRepository) AddItem(value string, userId int64) error {
	_, err := r.db.Exec(`
		INSERT INTO Items(value, numberInList, userId) VALUES(
			$1,
			(SELECT COUNT(*) FROM Items WHERE userId = $2) + 1,
			$2)`,
		value, userId)
	return err
}

// GetList возвращает список дел пользователя с указанным id.
func (r *itemsRepository) GetList(userId int64) ([]models.ListItem, error) {
	rows, err := r.db.Query("SELECT value, numberInList FROM Items WHERE userId = $1", userId)
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
