package dto

type OrderItemRequest struct {
	ProductID int `json:"product_id"`
	Qty       int     `json:"qty"`
    UnitPrice float64 `json:"unit_price"`
}
type PaymentRequest struct {
    Amount   float64 `json:"amount"`
    Method   string  `json:"method"`
    Currency string  `json:"currency"`
	Provider string `json:"provider"`
	Status string `json:"provider"`
}

type OrderRequest struct {
	UserID int `json:"userId"`
	Items []OrderItemRequest `json:"items"`
	Status string `json:"status"`
	TotalPrice float64 `json:"total_price"`
	Currency string `json:"currency"`
	Payments []PaymentRequest `json:"payments"`
}