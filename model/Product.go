package model

import (
	"database/sql"
	"strings"
	"time"
)

var (
	// DefaultProductBy ...
	DefaultProductBy = "def.updated_at"
	// ProductrBy ...
	ProductrBy = []string{
		"def.created_at", "def.updated_at",
	}

	// TypeDataMinus ...
	TypeDataMinus = "minus"
	// TypeDataPlus ...
	TypeDataPlus = "plus"

	productSelectString = `SELECT def.id, def.category_id, c.name, def.name, def.description, def.image_path, def.price, def.qty, def.created_at, def.updated_at, def.deleted_at FROM products def left join categories c on c.id = def.category_id `
)

func (model productModel) scanRows(rows *sql.Rows) (d ProductEntity, err error) {
	err = rows.Scan(
		&d.ID, &d.CategoryID, &d.CategoryName, &d.Name, &d.Description, &d.ImagePath, &d.Price, &d.Qty, &d.CreatedAt,
		&d.UpdatedAt, &d.DeletedAt,
	)

	return d, err
}

func (model productModel) scanRow(row *sql.Row) (d ProductEntity, err error) {
	err = row.Scan(
		&d.ID, &d.CategoryID, &d.CategoryName, &d.Name, &d.Description, &d.ImagePath, &d.Price, &d.Qty, &d.CreatedAt,
		&d.UpdatedAt, &d.DeletedAt,
	)

	return d, err
}

// productModel ...
type productModel struct {
	DB *sql.DB
}

// IProduct ...
type IProduct interface {
	FindAll(search, categoryID string, offset, limit int, by, sort string) ([]ProductEntity, int, error)
	FindByID(id string) (ProductEntity, error)
	UpdateStock(id string, qty int64, changedAt time.Time) error
}

// ProductEntity ....
type ProductEntity struct {
	ID           string         `db:"id"`
	CategoryID   string         `db:"category_id"`
	CategoryName string         `db:"category_name"`
	Name         string         `db:"name"`
	Description  string         `db:"description"`
	ImagePath    sql.NullString `db:"image_path"`
	Price        float64        `db:"price"`
	Qty          int            `db:"qty"`
	CreatedAt    string         `db:"created_at"`
	UpdatedAt    string         `db:"updated_at"`
	DeletedAt    sql.NullString `db:"deleted_at"`
}

// NewProductModel ...
func NewProductModel(db *sql.DB) IProduct {
	return &productModel{DB: db}
}

// FindAll ...
func (model productModel) FindAll(search, categoryID string, offset, limit int, by, sort string) (res []ProductEntity, count int, err error) {
	appendQuery := ``

	if categoryID != "" {
		appendQuery = ` AND def.category_id = '` + categoryID + `'`
	}

	query := productSelectString + ` WHERE def.deleted_at IS NULL ` + appendQuery + ` AND (LOWER ( def.name) LIKE ? ) ORDER BY ` + by + ` ` + sort + ` LIMIT ? OFFSET ? `
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

	query = `SELECT COUNT(def.id) FROM products def WHERE def.deleted_at IS NULL AND (LOWER ( def.name ) like ? )` + appendQuery
	err = model.DB.QueryRow(query, `%`+strings.ToLower(search)+`%`).Scan(&count)

	return res, count, err
}

// FindByID ...
func (model productModel) FindByID(id string) (res ProductEntity, err error) {
	query := productSelectString + ` WHERE def.deleted_at IS NULL AND def.id = ?
		ORDER BY def.created_at DESC LIMIT 1`
	row := model.DB.QueryRow(query, id)
	res, err = model.scanRow(row)

	return res, err
}

func (model productModel) UpdateStock(id string, qty int64, changedAt time.Time) (err error) {
	sql := `UPDATE products SET qty = ?, updated_at = ? WHERE deleted_at IS NULL
		AND id = ?`
	_, err = model.DB.Exec(sql, qty, changedAt, id)

	return err
}
