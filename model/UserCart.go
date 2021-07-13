package model

import (
	"database/sql"
	"retel-backend/usecase/viewmodel"
	"time"
)

var (
	// DefaultUserCartBy ...
	DefaultUserCartBy = "def.updated_at"
	// UserCartBy ...
	UserCartBy = []string{
		"def.created_at", "def.updated_at",
	}

	userCartSelectString = `SELECT def.id, def.user_id, def.product_id, def.qty, def.price, def.created_at, def.updated_at, def.deleted_at FROM user_cart def `
)

func (model userCartModel) scanRows(rows *sql.Rows) (d UserCartEntity, err error) {
	err = rows.Scan(
		&d.ID, &d.UserID, &d.ProductID, &d.Qty, &d.Price, &d.CreatedAt, &d.UpdatedAt, &d.DeletedAt,
	)

	return d, err
}

func (model userCartModel) scanRow(row *sql.Row) (d UserCartEntity, err error) {
	err = row.Scan(
		&d.ID, &d.UserID, &d.ProductID, &d.Qty, &d.Price, &d.CreatedAt, &d.UpdatedAt, &d.DeletedAt,
	)

	return d, err
}

// userCartModel ...
type userCartModel struct {
	DB *sql.DB
}

// IUserCart ...
type IUserCart interface {
	FindAll(userID string, offset, limit int, by, sort string) ([]UserCartEntity, int, error)
	Store(id string, body viewmodel.UserCartVM, changedAt time.Time) error
	Update(id string, body viewmodel.UserCartVM, changedAt time.Time) error
	Destroy(id string, changedAt time.Time) error
}

// UserCartEntity ....
type UserCartEntity struct {
	ID        string         `db:"id"`
	UserID    string         `db:"user_id"`
	ProductID string         `db:"product_id"`
	Qty       int64          `db:"qty"`
	Price     float64        `db:"price"`
	CreatedAt string         `db:"created_at"`
	UpdatedAt string         `db:"updated_at"`
	DeletedAt sql.NullString `db:"deleted_at"`
}

// NewUserCartModel ...
func NewUserCartModel(db *sql.DB) IUserCart {
	return &userCartModel{DB: db}
}

// FindAll ...
func (model userCartModel) FindAll(userID string, offset, limit int, by, sort string) (res []UserCartEntity, count int, err error) {
	query := userCartSelectString + ` WHERE def.deleted_at IS NULL AND def.user_id = ? ORDER BY ` + by + ` ` + sort + ` LIMIT ? OFFSET ? `
	rows, err := model.DB.Query(query, userID, limit, offset)
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

	query = `SELECT COUNT(def.id) FROM user_cart def WHERE def.deleted_at IS NULL AND def.user_id = ?`
	err = model.DB.QueryRow(query, userID).Scan(&count)

	return res, count, err
}

// Store ...
func (model userCartModel) Store(id string, body viewmodel.UserCartVM, changedAt time.Time) (err error) {
	sql := `INSERT INTO user_cart (id, user_id, product_id, qty, price, created_at, updated_at
		) VALUES(?, ?, ?, ?, ?, ?, ?)`
	_, err = model.DB.Exec(sql, id, body.UserID, body.ProductID, body.Qty, body.Price, changedAt, changedAt)

	return err
}

// Update ...
func (model userCartModel) Update(id string, body viewmodel.UserCartVM, changedAt time.Time) (err error) {
	sql := `UPDATE user_cart SET qty = ?, updated_at = ? WHERE deleted_at IS NULL
		AND id = ?`
	_, err = model.DB.Exec(sql, body.Qty, changedAt, id)

	return err
}

// Destroy ...
func (model userCartModel) Destroy(id string, changedAt time.Time) (err error) {
	sql := `UPDATE user_cart SET updated_at = ?, deleted_at = ?
		WHERE deleted_at IS NULL AND id = ?`
	_, err = model.DB.Exec(sql, changedAt, changedAt, id)

	return err
}
