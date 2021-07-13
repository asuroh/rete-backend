package request

// UserCartRequest ...
type UserCartRequest struct {
	ProductID string `json:"product_id" validate:"required"`
	Qty       int64  `json:"qty" validate:"required"`
}

// UserCartUpdateRequest ...
type UserCartUpdateRequest struct {
	Qty int64 `json:"qty" validate:"required"`
}
