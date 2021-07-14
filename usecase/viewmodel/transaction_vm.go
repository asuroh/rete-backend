package viewmodel

// TransactionVM ...
type TransactionVM struct {
	ID         string  `json:"id"`
	UserID     string  `json:"user_id"`
	Total      float64 `json:"total"`
	UrlPayment string  `json:"url_payment"`
	InvoceID   string  `json:"invoce_id"`
	Note       string  `json:"note"`
	Status     string  `json:"status"`
	Code       string  `json:"code"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
	DeletedAt  string  `json:"deleted_at"`
}
