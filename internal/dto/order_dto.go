package dto

type OrderItemDTO struct {
    ID        int     `json:"id"`
    ProductID int     `json:"productId"`
    Qty       int     `json:"qty"`
    UnitPrice float64 `json:"unitPrice"`
}

type PaymentDTO struct {
    ID         int     `json:"id"`
    Provider   string  `json:"provider"`
    AmountCents float64 `json:"amountCents"`
    Method     string  `json:"method"`
    Currency   string  `json:"currency"`
    Status     string  `json:"status"`
}

type OrderDTO struct {
    ID         int            `json:"id"`
    UserID     int            `json:"userId"`
    Status     string         `json:"status"`
    TotalPrice float64        `json:"total_price"`
    Currency   string         `json:"currency"`
    Items      []OrderItemDTO `json:"items"`
    Payments   []PaymentDTO   `json:"payments"`
    CreatedAt  string         `json:"created_at"`
    UpdatedAt  string         `json:"updated_at"`
}