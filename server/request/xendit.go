package request

//XenditRequest ...
type XenditRequest struct {
	TransactionID string `json:"transaction_id" validate:"required"`
	TokenID       string `json:"token_id" validate:"required"`
	ExternalID    string `json:"external_id" validate:"required"`
	Amount        int    `json:"amount" validate:"required"`
	CardCVN       string `json:"card_cvn" validate:"required"`
}

//XenditRefundRequest ...
type XenditRefundRequest struct {
	TransactionID string `json:"transaction_id" validate:"required"`
	ChargeID      string `json:"charge_id" validate:"required"`
	ExternalID    string `json:"external_id" validate:"required"`
	Amount        int    `json:"amount" validate:"required"`
}

// XenditInvoiceRequest ...
type XenditInvoiceRequest struct {
	ExternalID         string  `json:"external_id"`
	PayerEmail         string  `json:"payer_email"`
	Description        string  `json:"description"`
	ShouldSendEmail    bool    `json:"should_send_email"`
	Amount             float64 `json:"amount"`
	SuccessRedirectUrl string  `json:"success_redirect_url"`
	FailureRedirectUrl string  `json:"failure_redirect_url"`
}

// XenditInvoiceCallbackRequest ...
type XenditInvoiceCallbackRequest struct {
	ID                     string `json:"id"`
	UserID                 string `json:"user_id"`
	ExternalID             string `json:"external_id"`
	IsHigh                 bool   `json:"is_high"`
	MerchantName           string `json:"merchant_name"`
	Amount                 int64  `json:"amount"`
	Status                 string `json:"status"`
	PayerEmail             string `json:"payer_email"`
	Description            string `json:"description"`
	FeesPaidAmount         int64  `json:"fees_paid_amount"`
	AdjustedReceivedAmount int64  `json:"adjusted_received_amount"`
	BankCode               string `json:"bank_code"`
	RetailOutletName       string `json:"retail_outlet_name"`
	EwalletType            string `json:"ewallet_type"`
	OnDemandLink           string `json:"on_demand_link"`
	RecurringPaymentID     string `json:"recurring_payment_id"`
	PaidAmount             int64  `json:"paid_amount"`
	Updated                string `json:"updated"`
	Created                string `json:"created"`
	MidLabel               string `json:"mid_label"`
	Currency               string `json:"currency"`
	SuccessRedirectUrl     string `json:"success_redirect_url"`
	FailureRedirectUrl     string `json:"failure_redirect_url"`
	PaidAt                 string `json:"paid_at"`
	CreditCardChargeID     string `json:"credit_card_charge_id"`
	PaymentMethod          string `json:"payment_method"`
	PaymentChannel         string `json:"payment_channel"`
	PaymentDestination     string `json:"payment_destination"`
}
