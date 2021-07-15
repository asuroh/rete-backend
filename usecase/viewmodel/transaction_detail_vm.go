package viewmodel

// TransactionDetailVM ...
type TransactionDetailVM struct {
	ID            string  `json:"id"`
	TransactionID string  `json:"transaction_id"`
	ProductName   string  `json:"product_name"`
	ProductID     string  `json:"product_id"`
	Price         float64 `json:"price"`
	Qty           int64   `json:"qty"`
	Total         float64 `json:"total"`
	CheckIn       string  `json:"check_in"`
	CheckOut      string  `json:"check_out"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
	DeletedAt     string  `json:"deleted_at"`
}
