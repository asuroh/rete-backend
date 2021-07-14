package model

import (
	"database/sql"
	"retel-backend/usecase/viewmodel"
	"time"
)

var (
	// DefaultTransactionBy ...
	DefaultTransactionBy = "def.updated_at"
	// TransactionBy ...
	TransactionBy = []string{
		"def.created_at", "def.updated_at",
	}

	transactionSelectString = `SELECT def.id, def.user_id, def.total, def.url_payment, def.invoice_id, def.note, def.status, def.code, def.created_at, def.updated_at, def.deleted_at, usr.name FROM transaction def left join users usr on def.user_id = usr.id and usr.deleted_at is null`
)

func (model transactionModel) scanRows(rows *sql.Rows) (d TransactionEntity, err error) {
	err = rows.Scan(
		&d.ID, &d.UserID, &d.Total, &d.UrlPayment, &d.InvoceID, &d.Note, &d.Status, &d.Code, &d.CreatedAt, &d.UpdatedAt, &d.DeletedAt, &d.User.Name,
	)

	return d, err
}

func (model transactionModel) scanRow(row *sql.Row) (d TransactionEntity, err error) {
	err = row.Scan(
		&d.ID, &d.UserID, &d.Total, &d.UrlPayment, &d.InvoceID, &d.Note, &d.Status, &d.Code, &d.CreatedAt, &d.UpdatedAt, &d.DeletedAt, &d.User.Name,
	)

	return d, err
}

// transactionModel ...
type transactionModel struct {
	DB *sql.DB
}

// Itransaction ...
type Itransaction interface {
	FindAll(userID string, offset, limit int, by, sort string) ([]TransactionEntity, int, error)
	FindByID(id string) (TransactionEntity, error)
	Store(id string, body viewmodel.TransactionVM, changedAt time.Time) error
}

// TransactionEntity ....
type TransactionEntity struct {
	ID         string         `db:"id"`
	UserID     string         `db:"user_id"`
	User       UserEntity     `db:"user"`
	Total      float64        `db:"total"`
	UrlPayment sql.NullString `db:"url_payment"`
	InvoceID   sql.NullString `db:"invoce_id"`
	Note       sql.NullString `db:"note"`
	Status     string         `db:"status"`
	Code       string         `db:"code"`
	CreatedAt  string         `db:"created_at"`
	UpdatedAt  string         `db:"updated_at"`
	DeletedAt  sql.NullString `db:"deleted_at"`
}

// NewTransactionModel ...
func NewTransactionModel(db *sql.DB) Itransaction {
	return &transactionModel{DB: db}
}

// FindAll ...
func (model transactionModel) FindAll(userID string, offset, limit int, by, sort string) (res []TransactionEntity, count int, err error) {
	appendQuery := ``
	if userID != "" {
		appendQuery = ` AND def.user_id = '` + userID + `' `
	}
	query := transactionSelectString + ` WHERE def.deleted_at IS NULL ` + appendQuery + ` ORDER BY 
	` + by + ` ` + sort + ` LIMIT ? OFFSET ? `
	rows, err := model.DB.Query(query, limit, offset)
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

	query = `SELECT COUNT(def.id) FROM transaction def WHERE def.deleted_at IS NULL ` + appendQuery
	err = model.DB.QueryRow(query).Scan(&count)

	return res, count, err
}

// FindByID ...
func (model transactionModel) FindByID(id string) (res TransactionEntity, err error) {
	query := transactionSelectString + ` WHERE def.deleted_at IS NULL AND def.id = ?
		ORDER BY def.created_at DESC LIMIT 1`
	row := model.DB.QueryRow(query, id)
	res, err = model.scanRow(row)

	return res, err
}

// Store ...
func (model transactionModel) Store(id string, body viewmodel.TransactionVM, changedAt time.Time) (err error) {
	sql := `INSERT INTO transaction (id, user_id, url_payment, invoce_id, note, status, created_at, updated_at
		) VALUES(?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = model.DB.Exec(sql, id, body.UserID, body.UrlPayment, body.InvoceID, body.Note, body.Status, changedAt, changedAt)

	return err
}
