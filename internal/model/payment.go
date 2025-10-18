package model

import (
	"errors"
	"time"
)

type Payment struct {
    ID                    int      `gorm:"primaryKey" json:"id"`
    OrderID               int      `json:"order_id"`
    Order                 Order    `gorm:"foreignKey:OrderID" json:"-"`
    Provider              string   `json:"provider"` // e.g. "stripe"
    AmountCents           int64    `json:"amount_cents"` // Changed to int64 for Stripe compatibility
    Method                string   `json:"method"`
    Currency              string   `json:"currency"`
    Status                string   `json:"status"`   // pending, succeeded, failed, canceled
    StripePaymentIntentID *string  `json:"stripe_payment_intent_id,omitempty"`
    StripeChargeID        *string  `json:"stripe_charge_id,omitempty"`
    StripeCustomerID      *string  `json:"stripe_customer_id,omitempty"`
    FailureReason         *string  `json:"failure_reason,omitempty"`
    FailureCode           *string  `json:"failure_code,omitempty"`
    ProcessedAt           *time.Time `json:"processed_at,omitempty"`
    CreatedAt             time.Time `json:"created_at"`
    UpdatedAt             time.Time `json:"updated_at"`
}

func NewPayment(orderID int, provider, method, currency, status string, amountCents int64) (*Payment, error) {
    now := time.Now()
    payment := &Payment{
        OrderID:     orderID,
        Provider:    provider,
        Method:      method,
        Currency:    currency,
        Status:      status,
        AmountCents: amountCents,
        CreatedAt:   now,
        UpdatedAt:   now,
    }
    err := payment.validatePayment()
    if err != nil {
        return nil, err
    }
    return payment, nil
}

func (p *Payment) validatePayment() error {
    if p.OrderID == 0 {
        return errors.New("order ID is required")
    }
    if p.Provider == "" {
        return errors.New("provider is required") 
    }
    if p.Method == "" {
        return errors.New("method is required")
    }
    if p.Currency == "" {
        return errors.New("currency is required")
    }
    if p.Status == "" {
        return errors.New("status is required")
    }
    if p.AmountCents <= 0 { 
        return errors.New("amount cents must be positive")
    }
    
    // Additional validations for better Stripe compatibility
    validStatuses := []string{"pending", "succeeded", "failed", "canceled", "requires_payment_method", "requires_confirmation"}
    validStatus := false
    for _, status := range validStatuses {
        if p.Status == status {
            validStatus = true
            break
        }
    }
    if !validStatus {
        return errors.New("invalid payment status")
    }
    
    if len(p.Currency) != 3 {
        return errors.New("currency must be a 3-letter ISO code")
    }
    
    return nil
}