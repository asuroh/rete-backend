package request

// TransactionDetailRequest ...
type TransactionDetailRequest struct {
	TransactionID string `json:"transaction_id"`
	ProductID     string `json:"product_id"`
	Qty           int64  `json:"qty"`
	CheckIn       string `json:"check_in"`
	CheckOut      string `json:"check_out"`
}
