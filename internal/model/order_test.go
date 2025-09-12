package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestNewOrder(t *testing.T) {
	items := []OrderItem{
		{
			ID:        1,
			OrderID:   10,
			ProductID: 1,
			Product: nil,
			Qty:       1,
			UnitPrice: 150.0,
		},
	}

	payments := []Payment{
		{
			ID: 1,
			OrderID: 1,
			Provider: "stripe",
			AmountCents: 150.0,
			Method: "credit_card",
			Currency: "USD",
			Status: "pending",
			StripePaymentIntentID: nil,
			StripeChargeID: nil,
		},
	}

	order, err := NewOrder(1, items, "pending", 150.0, "USD", payments)

	assert.NoError(t, err)
	assert.NotNil(t, order)
	
	// Assert basic order fields
	assert.Equal(t, 1, order.UserID)
	assert.Equal(t, "pending", order.Status)
	assert.Equal(t, 150.0, order.TotalPrice)
	assert.Equal(t, "USD", order.Currency)
	
	// Assert order has correct number of items and payments
	assert.Len(t, order.Items, 1)
	assert.Len(t, order.Payments, 1)
	
	// Assert item details
	assert.Equal(t, 1, order.Items[0].ProductID)
	assert.Equal(t, 1, order.Items[0].Qty)
	assert.Equal(t, 150.0, order.Items[0].UnitPrice)
	
	// Assert payment details
	assert.Equal(t, "stripe", order.Payments[0].Provider)
	assert.Equal(t, "credit_card", order.Payments[0].Method)
	assert.Equal(t, 150.0, order.Payments[0].AmountCents)
	assert.Equal(t, "USD", order.Payments[0].Currency)
	assert.Equal(t, "pending", order.Payments[0].Status)
}

func TestNewOrderWithMissingUserID(t *testing.T) {
	items := []OrderItem{
		{ProductID: 1, Qty: 1, UnitPrice: 150.0},
	}
	payments := []Payment{
		{Provider: "stripe", Method: "credit_card", Currency: "USD"},
	}

	order, err := NewOrder(0, items, "pending", 150.0, "USD", payments)
	
	assert.Error(t, err)
	assert.Nil(t, order)
	assert.Contains(t, err.Error(), "user ID is required")
}

func TestNewOrderWithNoItems(t *testing.T) {
	payments := []Payment{
		{Provider: "stripe", Method: "credit_card", Currency: "USD"},
	}

	order, err := NewOrder(1, []OrderItem{}, "pending", 0.0, "USD", payments)
	
	assert.Error(t, err)
	assert.Nil(t, order)
	assert.Contains(t, err.Error(), "order must have at least one item")
}

func TestNewOrderWithNegativeQuantity(t *testing.T) {
	items := []OrderItem{
		{ProductID: 1, Qty: -1, UnitPrice: 150.0},
	}
	payments := []Payment{
		{Provider: "stripe", Method: "credit_card", Currency: "USD"},
	}

	order, err := NewOrder(1, items, "pending", 150.0, "USD", payments)
	
	assert.Error(t, err)
	assert.Nil(t, order)
	assert.Contains(t, err.Error(), "item quantity must be positive")
}

func TestNewOrderWithZeroQuantity(t *testing.T) {
	items := []OrderItem{
		{ProductID: 1, Qty: 0, UnitPrice: 150.0},
	}
	payments := []Payment{
		{Provider: "stripe", Method: "credit_card", Currency: "USD"},
	}

	order, err := NewOrder(1, items, "pending", 150.0, "USD", payments)
	
	assert.Error(t, err)
	assert.Nil(t, order)
	assert.Contains(t, err.Error(), "item quantity must be positive")
}

func TestNewOrderTableDriven(t *testing.T) {
	validItems := []OrderItem{
		{ProductID: 1, Qty: 1, UnitPrice: 150.0},
	}
	validPayments := []Payment{
		{Provider: "stripe", Method: "credit_card", Currency: "USD"},
	}

	tests := []struct {
		name       string
		userID     int
		items      []OrderItem
		status     string
		totalPrice float64
		currency   string
		payments   []Payment
		wantErr    bool
	}{
		{
			name:       "valid order",
			userID:     1,
			items:      validItems,
			status:     "pending",
			totalPrice: 150.0,
			currency:   "USD",
			payments:   validPayments,
			wantErr:    false,
		},
		{
			name:       "missing user ID",
			userID:     0,
			items:      validItems,
			status:     "pending",
			totalPrice: 150.0,
			currency:   "USD",
			payments:   validPayments,
			wantErr:    true,
		},
		{
			name:		"quantity must be positive",
			userID: 1,
			items: []OrderItem{
				{ProductID: 1, Qty: -1, UnitPrice: 150.0},
			},
			status: "pending",
			totalPrice: 150.0,
			currency: "USD",
			payments: validPayments,
			wantErr: true,
		},
		{
			name:		"order must have at least one item",
			userID: 1,
			items: []OrderItem{},
			status: "pending",
			totalPrice: 150.0,
			currency: "USD",
			payments: validPayments,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order, err := NewOrder(tt.userID, tt.items, tt.status, tt.totalPrice, tt.currency, tt.payments)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, order)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, order)
				assert.Equal(t, tt.userID, order.UserID)
				assert.Equal(t, tt.status, order.Status)
				assert.Equal(t, tt.totalPrice, order.TotalPrice)
				assert.Equal(t, tt.currency, order.Currency)
				assert.Len(t, order.Items, len(tt.items))
				assert.Len(t, order.Payments, len(tt.payments))
			}
		})
	}
}