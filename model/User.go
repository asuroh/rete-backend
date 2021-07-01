package model

import (
	"database/sql"
	"fmt"
	"retel-backend/helper"
	"retel-backend/usecase/viewmodel"
	"strings"
	"time"
)

var (
	// DefaultUserBy ...
	DefaultUserBy = "def.updated_at"
	// UserBy ...
	UserBy = []string{
		"def.created_at", "def.updated_at",
	}

	userSelectString = `SELECT def.id, def.name, def.email, def.password, def.created_at, def.updated_at, def.deleted_at, def.image_path FROM users def`
)

func (model userModel) scanRows(rows *sql.Rows) (d UserEntity, err error) {
	err = rows.Scan(
		&d.ID, &d.Name, &d.Email, &d.Password, &d.CreatedAt,
		&d.UpdatedAt, &d.DeletedAt, &d.ImagePath,
	)

	return d, err
}

func (model userModel) scanRow(row *sql.Row) (d UserEntity, err error) {
	err = row.Scan(
		&d.ID, &d.Name, &d.Email, &d.Password, &d.CreatedAt,
		&d.UpdatedAt, &d.DeletedAt, &d.ImagePath,
	)

	return d, err
}

// userModel ...
type userModel struct {
	DB *sql.DB
}

// IUser ...
type IUser interface {
	FindAll(search string, offset, limit int, by, sort string) ([]UserEntity, int, error)
	FindByID(id string) (UserEntity, error)
	FindByEmail(email string) (UserEntity, error)
	Store(id string, body viewmodel.UserVM, changedAt time.Time) error
	Update(id string, body viewmodel.UserVM, changedAt time.Time) error
	UpdateImage(id, imagepath string, changedAt time.Time) error
	Destroy(id string, changedAt time.Time) error
}

// UserEntity ....
type UserEntity struct {
	ID        string         `db:"id"`
	Name      sql.NullString `db:"name"`
	Email     string         `db:"email"`
	ImagePath sql.NullString `db:"image_path"`
	Password  string         `db:"password"`
	CreatedAt string         `db:"created_at"`
	UpdatedAt string         `db:"updated_at"`
	DeletedAt sql.NullString `db:"deleted_at"`
}

// NewUserModel ...
func NewUserModel(db *sql.DB) IUser {
	return &userModel{DB: db}
}

// FindAll ...
func (model userModel) FindAll(search string, offset, limit int, by, sort string) (res []UserEntity, count int, err error) {
	query := userSelectString + ` WHERE def.deleted_at IS NULL AND (
	LOWER ( def.name) LIKE ? ) ORDER BY ` + by + ` ` + sort + ` LIMIT ? OFFSET ? `
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

	query = `SELECT COUNT(def.id) FROM users def WHERE def.deleted_at IS NULL AND (LOWER ( def.name ) like ? )`
	err = model.DB.QueryRow(query, `%`+strings.ToLower(search)+`%`).Scan(&count)

	return res, count, err
}

// FindByID ...
func (model userModel) FindByID(id string) (res UserEntity, err error) {
	query := userSelectString + ` WHERE def.deleted_at IS NULL AND def.id = ?
		ORDER BY def.created_at DESC LIMIT 1`
	row := model.DB.QueryRow(query, id)
	res, err = model.scanRow(row)

	return res, err
}

// FindByEmail ...
func (model userModel) FindByEmail(email string) (res UserEntity, err error) {
	query := userSelectString + ` WHERE def.deleted_at IS NULL  AND LOWER (def.email) = ? ORDER BY def.created_at DESC  LIMIT 1`
	row := model.DB.QueryRow(query, strings.ToLower(email))
	res, err = model.scanRow(row)
	if err != nil && err.Error() == helper.SQLHandlerErrorRowNull {
		err = nil
	}

	return res, err
}

// Store ...
func (model userModel) Store(id string, body viewmodel.UserVM, changedAt time.Time) (err error) {
	sql := `INSERT INTO users (id, name, email, password, created_at, updated_at
		) VALUES(?, ?, ?, ?, ?, ?)`
	_, err = model.DB.Exec(sql, id, body.Name, body.Email, body.Password, changedAt, changedAt)

	return err
}

// Update ...
func (model userModel) Update(id string, body viewmodel.UserVM, changedAt time.Time) (err error) {
	sql := `UPDATE users SET name = ?, email = ?, updated_at = ? WHERE deleted_at IS NULL
		AND "id" = ?`
	fmt.Println(body.Name)
	_, err = model.DB.Exec(sql, body.Name, body.Email, changedAt, id)

	return err
}

// UpdateImage ...
func (model userModel) UpdateImage(id, ImagePath string, changedAt time.Time) (err error) {
	sql := `UPDATE users SET image_path = ?, updated_at = ? WHERE deleted_at IS NULL
		AND id = ?`
	_, err = model.DB.Exec(sql, ImagePath, changedAt, id)

	return err
}

// Destroy ...
func (model userModel) Destroy(id string, changedAt time.Time) (err error) {
	sql := `UPDATE users SET updated_at = ?, deleted_at = ?
		WHERE deleted_at IS NULL AND id = ?`
	_, err = model.DB.Exec(sql, changedAt, changedAt, id)

	return err
}
