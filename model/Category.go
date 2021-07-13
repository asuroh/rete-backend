package model

import (
	"database/sql"
	"strings"
)

var (
	// DefaultProductBy ...
	DefaultCategoryBy = "def.updated_at"
	// ProductrBy ...
	CategoryBy = []string{
		"def.created_at", "def.updated_at",
	}

	categorySelectString = `SELECT def.id, def.name, def.description, def.created_at, def.updated_at, def.deleted_at FROM categories def `
)

func (model categoryModel) scanRows(rows *sql.Rows) (d CategoryEntity, err error) {
	err = rows.Scan(
		&d.ID, &d.Name, &d.Description, &d.CreatedAt, &d.UpdatedAt, &d.DeletedAt,
	)

	return d, err
}

func (model categoryModel) scanRow(row *sql.Row) (d CategoryEntity, err error) {
	err = row.Scan(
		&d.ID, &d.Name, &d.Description, &d.CreatedAt, &d.UpdatedAt, &d.DeletedAt,
	)

	return d, err
}

// categoryModel ...
type categoryModel struct {
	DB *sql.DB
}

// ICategory ...
type ICategory interface {
	FindAll(search string, offset, limit int, by, sort string) ([]CategoryEntity, int, error)
}

// CategoryEntity ....
type CategoryEntity struct {
	ID          string         `db:"id"`
	Name        string         `db:"name"`
	Description string         `db:"description"`
	CreatedAt   string         `db:"created_at"`
	UpdatedAt   string         `db:"updated_at"`
	DeletedAt   sql.NullString `db:"deleted_at"`
}

// NewCategoryModel ...
func NewCategoryModel(db *sql.DB) ICategory {
	return &categoryModel{DB: db}
}

// FindAll ...
func (model categoryModel) FindAll(search string, offset, limit int, by, sort string) (res []CategoryEntity, count int, err error) {
	query := categorySelectString + ` WHERE def.deleted_at IS NULL AND (LOWER ( def.name) LIKE ? ) ORDER BY ` + by + ` ` + sort + ` LIMIT ? OFFSET ? `
	rows, err := model.DB.Query(query, `%`+strings.ToLower(search)+`%`, limit, offset)
	if err != nil {
		return res, count, err
	}
	defer rows.Close()

	for rows.Next() {
		d, err := model.scanRows(rows)
		if err != nil {
			return res, count, err
		}
		res = append(res, d)
	}
	err = rows.Err()
	if err != nil {
		return res, count, err
	}

	query = `SELECT COUNT(def.id) FROM products def WHERE def.deleted_at IS NULL AND (LOWER ( def.name ) like ? )`
	err = model.DB.QueryRow(query, `%`+strings.ToLower(search)+`%`).Scan(&count)

	return res, count, err
}
