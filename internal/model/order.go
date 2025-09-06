package model

import (
	"errors"
	"time"
)

type Order struct {
	ID int `json:"id" gorm:"primaryKey"`
	UserID int `json:"user_id"`
	User User `gorm:"foreignKey:UserID" json:"-"`
	Status string `json:"status"`
	TotalPrice float64 `json:"total_price"`
	Currency string `json:"currency"`
	Items []OrderItem `gorm:"foreignKey:OrderID" json:"items"`
	Payments []Payment `gorm:"foreignKey:OrderID" json:"payments"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewOrder(
    userID int,
    items []OrderItem,
    status string,
    totalPrice float64,
    currency string,
    payments []Payment,
) (*Order, error) {
    order := &Order{
        UserID:     userID,
        Items:      items,
        Status:     status,
        TotalPrice: totalPrice,
        Currency:   currency,
        Payments:   payments,
    }
    if err := order.Validate(); err != nil {
        return nil, err
    }
    return order, nil
}

func (o *Order) Validate() error {
	if o.UserID == 0 {
		return errors.New("user ID is required")
	}
	if len(o.Items) == 0 {
		return errors.New("order must have at least one item")
	}
	for _, item := range o.Items {
		if item.Qty <= 0 {
			return errors.New("item quantity must be positive")
		}
	}
	return nil
}