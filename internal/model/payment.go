package model

import "time"

type Payment struct {
    ID        int64      `gorm:"primaryKey" json:"id"`
    OrderID   int64      `json:"orderId"`
    Order     Order     `gorm:"foreignKey:OrderID"`
    Provider  string    `json:"provider"` // e.g. "stripe"
    AmountCents int     `json:"amountCents"`
    Currency  string    `json:"currency"`
    Status    string    `json:"status"`   // pending, succeeded, failed
    StripePaymentIntentID *string `json:"stripePaymentIntentId,omitempty"`
    StripeChargeID        *string `json:"stripeChargeId,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
