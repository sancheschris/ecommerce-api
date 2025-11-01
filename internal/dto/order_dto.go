package dto

type ProductDTO struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Active    bool    `json:"active"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type OrderItemDTO struct {
	ID        int        `json:"id"`
	OrderID   int        `json:"order_id"`
	ProductID int        `json:"product_id"`
	Product   ProductDTO `json:"product"`
	Qty       int        `json:"qty"`
	UnitPrice float64    `json:"unit_price"`
}

type PaymentDTO struct {
	ID        int     `json:"id"`
	OrderID   int     `json:"order_id"`
	Provider  string  `json:"provider"`
	Amount    float64 `json:"amount"`
	Method    string  `json:"method"`
	Currency  string  `json:"currency"`
	Status    string  `json:"status"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type CreateOrderItemRequest struct {
    ProductID int     `json:"product_id" validate:"required,min=1"`
    Quantity  int     `json:"quantity" validate:"required,min=1"`
    Price     float64 `json:"price" validate:"required,min=0"`
}

type CreateOrderRequest struct {
    UserID int                      `json:"user_id" validate:"required,min=1"`
    Items  []CreateOrderItemRequest `json:"items" validate:"required,min=1"`
}

type UpdateOrderItemRequest struct {
    ID        int     `json:"id,omitempty"`        // Optional for new items
    ProductID int     `json:"product_id" validate:"required,min=1"`
    Quantity  int     `json:"quantity" validate:"required,min=1"`
    Price     float64 `json:"price" validate:"required,min=0"`
}

type UpdateOrderRequest struct {
    Status string                   `json:"status,omitempty"`    
    Items  []UpdateOrderItemRequest `json:"items,omitempty"`      
}

type OrderDTO struct {
	ID         int            `json:"id"`
	UserID     int            `json:"user_id"`
	Status     string         `json:"status"`
	TotalPrice float64        `json:"total_price"`
	Currency   string         `json:"currency"`
	Items      []OrderItemDTO `json:"items"`
	Payments   []PaymentDTO   `json:"payments"`
	CreatedAt  string         `json:"created_at"`
	UpdatedAt  string         `json:"updated_at"`
}
