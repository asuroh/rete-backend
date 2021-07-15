package model

import (
	"database/sql"
	"retel-backend/usecase/viewmodel"
	"time"
)

var (
	// DefaultTransactioDetailnBy ...
	DefaultTransactioDetailnBy = "def.updated_at"
	// TransactionDetailBy ...
	TransactionDetailBy = []string{
		"def.created_at", "def.updated_at",
	}

	transactionDetailSelectString = `SELECT def.id, def.user_id, def.total, def.url_payment, def.invoice_id, def.note, def.created_at, def.updated_at, def.deleted_at FROM transaction def `
)

func (model transactionDetailModel) scanRows(rows *sql.Rows) (d TransactionDetailEntity, err error) {
	err = rows.Scan(
		&d.ID, &d.TransactionID, &d.ProductName, &d.ProductID, &d.Price, &d.Qty, &d.CreatedAt, &d.UpdatedAt, &d.DeletedAt,
	)

	return d, err
}

func (model transactionDetailModel) scanRow(row *sql.Row) (d TransactionDetailEntity, err error) {
	err = row.Scan(
		&d.ID, &d.TransactionID, &d.ProductName, &d.ProductID, &d.Price, &d.Qty, &d.CreatedAt, &d.UpdatedAt, &d.DeletedAt,
	)

	return d, err
}

// transactionDetailModel ...
type transactionDetailModel struct {
	DB *sql.DB
}

// ItransactionDetail ...
type ItransactionDetail interface {
	Store(id string, body viewmodel.TransactionDetailVM, changedAt time.Time) error
}

// TransactionDetailEntity ....
type TransactionDetailEntity struct {
	ID            string         `db:"id"`
	TransactionID string         `db:"transaction_id"`
	ProductName   string         `db:"product_name"`
	ProductID     string         `db:"product_id"`
	Price         float64        `db:"price"`
	Qty           int64          `db:"qty"`
	CheckIn       string         `db:"check_in"`
	CheckOut      string         `db:"check_out"`
	CreatedAt     string         `db:"created_at"`
	UpdatedAt     string         `db:"updated_at"`
	DeletedAt     sql.NullString `db:"deleted_at"`
}

// NewTransactionDetailModel ...
func NewTransactionDetailModel(db *sql.DB) ItransactionDetail {
	return &transactionDetailModel{DB: db}
}

// Store ...
func (model transactionDetailModel) Store(id string, body viewmodel.TransactionDetailVM, changedAt time.Time) (err error) {
	sql := `INSERT INTO transaction_details (id, transaction_id, product_name, product_id, price, qty, check_in, check_out, created_at, updated_at
		) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = model.DB.Exec(sql, id, body.TransactionID, body.ProductName, body.ProductID, body.Price, body.Qty, body.CheckIn, body.CheckOut, changedAt, changedAt)

	return err
}
