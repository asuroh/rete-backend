package request

// TransactionRequest ...
type TransactionRequest struct {
	UserID     string  `json:"user_id"`
	Total      float64 `json:"total"`
	UrlPayment string  `json:"url_payment"`
	InvoceID   string  `json:"invoce_id"`
	Note       string  `json:"note"`
}
