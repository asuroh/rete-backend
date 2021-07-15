package viewmodel

// InvoiceData ...
type InvoiceData struct {
	ID                        string                     `json:"id"`
	UserID                    string                     `json:"user_id"`
	ExternalID                string                     `json:"external_id"`
	Status                    string                     `json:"status"`
	MerchantName              string                     `json:"merchant_name"`
	MerchantProfilePictureUrl string                     `json:"merchant_profile_picture_url"`
	Amount                    int64                      `json:"amount"`
	PayerEmail                string                     `json:"payer_email"`
	Description               string                     `json:"description"`
	InvoiceUrl                string                     `json:"invoice_url"`
	ExpiryDate                string                     `json:"expiry_date"`
	AvailableBanks            []AvailableBanksVM         `json:"available_banks"`
	AvailableRetailOutlets    []AvailableRetailOutletsVM `json:"available_retail_outlets"`
	TransferAmount            int64                      `json:"transfer_amount"`
	ShouldExcludeCreditCard   bool                       `json:"should_exclude_credit_card"`
	ShouldSendEmail           bool                       `json:"should_send_email"`
	Created                   string                     `json:"created"`
	Updated                   string                     `json:"updated"`
	MidLabel                  string                     `json:"mid_label"`
	Currency                  string                     `json:"currency"`
	FixedVa                   bool                       `json:"fixed_va"`
}

// InvoiceError ...
type InvoiceError struct {
	ErrorCode string `json:"error_code"`
	Message   string `json:"message"`
}

// AvailableBanksVM ...
type AvailableBanksVM struct {
	BankCode          string `json:"bank_code"`
	CollectionType    string `json:"collection_type"`
	BankAccountNumber string `json:"bank_account_number"`
	TransferAmount    int64  `json:"transfer_amount"`
	BankBranch        string `json:"bank_branch"`
	AccountHolderName string `json:"account_holder_name"`
}

// AvailableRetailOutletsVM ...
type AvailableRetailOutletsVM struct {
	RetailOutletName string `json:"retail_outlet_name"`
	PaymentCode      string `json:"payment_code"`
	TransferAmount   int64  `json:"transfer_amount"`
}
