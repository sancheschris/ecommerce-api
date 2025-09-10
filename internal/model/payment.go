package model

import "time"

type Payment struct {
    ID        int      `gorm:"primaryKey" json:"id"`
    OrderID   int      `json:"order_id"`
    Order     Order     `gorm:"foreignKey:OrderID" json:"-"`
    Provider  string    `json:"provider"` // e.g. "stripe"
    AmountCents float64     `json:"amount_cents"`
    Method string `json:"method"`
    Currency  string    `json:"currency"`
    Status    string    `json:"status"`   // pending, succeeded, failed
    StripePaymentIntentID *string `json:"stripePaymentIntentId,omitempty"`
    StripeChargeID        *string `json:"stripeChargeId,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
