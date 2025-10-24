package dto

type CreatePaymentRequest struct {
	OrderID int `json:"order_id"`
	Amount int64 `json:"amount_cents"`
	Currency string `json:"currency"`
}

