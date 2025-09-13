package payment

import (
	"testing"

	"github.com/sancheschris/ecommerce-api/internal/model"
	"github.com/sancheschris/ecommerce-api/internal/repository"
	"github.com/stretchr/testify/assert"
)


func TestNewPayment(t *testing.T) {
	db := repository.SetupTestDB(model.Payment{})
	paymentDB := NewPayment(db)

	payment, err := model.NewPayment(1, "stripe", "credit_card", "USD", "pending", 10)
	if err != nil {
		t.Fatalf("Cannot create payment %s", err)
	}

	err = paymentDB.Create(payment)
	if err != nil {
		t.Error(err)
	}
}

func TestCreateNewPayment(t *testing.T) {
	db := repository.SetupTestDB(model.Payment{})
	paymentDB := NewPayment(db)

	tests := []struct {
		name        string
		orderID     int
		provider    string
		amountCents int64
		method      string
		currency    string
		status      string
		wantErr     bool
	}{
		{
			name:        "Valid Payment - USD",
			orderID:     1,
			provider:    "stripe",
			amountCents: 1500, // $15.00
			method:      "credit_card",
			currency:    "USD",
			status:      "pending",
			wantErr:     false,
		},
		{
			name:        "Valid Payment - EUR",
			orderID:     2,
			provider:    "stripe",
			amountCents: 2599, // €25.99
			method:      "credit_card",
			currency:    "EUR",
			status:      "pending",
			wantErr:     false,
		},
		{
			name:        "Valid Payment - Large Amount",
			orderID:     3,
			provider:    "stripe",
			amountCents: 100000, // $1000.00
			method:      "credit_card",
			currency:    "USD",
			status:      "succeeded",
			wantErr:     false,
		},
		{
			name:        "Invalid Payment - Negative Amount",
			orderID:     4,
			provider:    "stripe",
			amountCents: -100,
			method:      "credit_card",
			currency:    "USD",
			status:      "pending",
			wantErr:     true,
		},
		{
			name:        "Invalid Payment - Zero Amount",
			orderID:     5,
			provider:    "stripe",
			amountCents: 0,
			method:      "credit_card",
			currency:    "USD",
			status:      "pending",
			wantErr:     true,
		},
		{
			name:        "Invalid Payment - Empty Provider",
			orderID:     6,
			provider:    "",
			amountCents: 1000,
			method:      "credit_card",
			currency:    "USD",
			status:      "pending",
			wantErr:     true,
		},
		{
			name:        "Invalid Payment - Empty Method",
			orderID:     7,
			provider:    "stripe",
			amountCents: 1000,
			method:      "",
			currency:    "USD",
			status:      "pending",
			wantErr:     true,
		},
		{
			name:        "Invalid Payment - Invalid Currency (too short)",
			orderID:     8,
			provider:    "stripe",
			amountCents: 1000,
			method:      "credit_card",
			currency:    "US",
			status:      "pending",
			wantErr:     true,
		},
		{
			name:        "Invalid Payment - Invalid Currency (too long)",
			orderID:     9,
			provider:    "stripe",
			amountCents: 1000,
			method:      "credit_card",
			currency:    "USDX",
			status:      "pending",
			wantErr:     true,
		},
		{
			name:        "Invalid Payment - Invalid Status",
			orderID:     10,
			provider:    "stripe",
			amountCents: 1000,
			method:      "credit_card",
			currency:    "USD",
			status:      "invalid_status",
			wantErr:     true,
		},
		{
			name:        "Valid Payment - Different Provider",
			orderID:     11,
			provider:    "paypal",
			amountCents: 1500,
			method:      "paypal_account",
			currency:    "USD",
			status:      "pending",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payment, err := model.NewPayment(tt.orderID, tt.provider, tt.method, tt.currency, tt.status, tt.amountCents)
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, payment)
				return
			}
			
			assert.NoError(t, err)
			assert.NotNil(t, payment)
			
			// Test database creation
			err = paymentDB.Create(payment)
			assert.NoError(t, err)
			
			// Verify the payment was saved with correct data
			assert.NotZero(t, payment.ID, "Payment ID should be set after creation")
			
			// Retrieve and verify
			savedPayment, err := paymentDB.GetByID(payment.ID)
			assert.NoError(t, err)
			assert.Equal(t, tt.orderID, savedPayment.OrderID)
			assert.Equal(t, tt.provider, savedPayment.Provider)
			assert.Equal(t, tt.method, savedPayment.Method)
			assert.Equal(t, tt.currency, savedPayment.Currency)
			assert.Equal(t, tt.status, savedPayment.Status)
			assert.Equal(t, tt.amountCents, savedPayment.AmountCents)
		})
	}
}

func TestGetByID(t *testing.T) {
	db := repository.SetupTestDB(model.Payment{})
	paymentDB := NewPayment(db)

	payment, err := model.NewPayment(1, "stripe", "credit_card", "EUR", "succeeded", 1099)
	assert.NoError(t, err)
	err = paymentDB.Create(payment)
	assert.NoError(t, err)

	tests := []struct {
		name       string
		id         int
		wantErr    bool
		wantResult bool
	}{
		{
			name:       "Valid ID - Existing Payment",
			id:         payment.ID,
			wantErr:    false,
			wantResult: true,
		},
		{
			name:       "Invalid ID - Non-existent",
			id:         99999,
			wantErr:    true,
			wantResult: false,
		},
		{
			name:       "Invalid ID - Zero",
			id:         0,
			wantErr:    true,
			wantResult: false,
		},
		{
			name:       "Invalid ID - Negative",
			id:         -1,
			wantErr:    true,
			wantResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := paymentDB.GetByID(tt.id)
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				// Detailed validation for happy path
				assert.Equal(t, tt.id, result.ID)
				assert.Equal(t, "EUR", result.Currency)
				assert.Equal(t, "stripe", result.Provider)
				assert.Equal(t, int64(1099), result.AmountCents)
				assert.Equal(t, "succeeded", result.Status)
			}
		})
	}
}

func TestUpdate_ValidPayment_UpdatesFieldsSuccessfully(t *testing.T) {
	db := repository.SetupTestDB(model.Payment{})
	paymentDB := NewPayment(db)

	payment, err := model.NewPayment(1, "stripe", "credit_card", "EUR", "pending", 1099)
	assert.NoError(t, err)
	err = paymentDB.Create(payment)
	assert.NoError(t, err)

	existingPayment, err := paymentDB.GetByID(payment.ID)
	assert.NoError(t, err)

	// Store original values for comparison
	originalAmount := existingPayment.AmountCents
	originalCurrency := existingPayment.Currency
	originalStatus := existingPayment.Status

	// Modify the payment
	existingPayment.AmountCents = 2299
	existingPayment.Currency = "AUD"
	existingPayment.Status = "succeeded"

	err = paymentDB.Update(existingPayment)
	assert.NoError(t, err)

	updatedPayment, err := paymentDB.GetByID(payment.ID)
	assert.NoError(t, err)
	assert.NotNil(t, updatedPayment)

	// Verifiy the values changed
	assert.Equal(t, existingPayment.AmountCents, updatedPayment.AmountCents)
	assert.Equal(t, existingPayment.Currency, updatedPayment.Currency)
	assert.Equal(t, existingPayment.Status, updatedPayment.Status)

	// Verify they're different from original values
	assert.NotEqual(t, originalAmount, updatedPayment.AmountCents)
	assert.NotEqual(t, originalCurrency, updatedPayment.Currency)
	assert.NotEqual(t, originalStatus, updatedPayment.Status)
}

func TestUpdate_InvalidPayment_ReturnsError(t *testing.T) {
	db := repository.SetupTestDB(model.Payment{})
	paymentDB := NewPayment(db)

	payment, err := model.NewPayment(1, "stripe", "credit_card", "EUR", "pending", 1099)
	assert.NoError(t, err)

	payment.ID = 99999 // Non-existent ID

	err = paymentDB.Update(payment)
	assert.Error(t, err) 
}

func TestDelete_ValidPayment_DeletesSuccessfully(t *testing.T) {
	db := repository.SetupTestDB(model.Payment{})
	paymentDB := NewPayment(db)

	payment, err := model.NewPayment(1, "stripe", "credit_card", "CAD", "pending", 3189)
	assert.NoError(t, err)

	err = paymentDB.Create(payment)
	assert.NoError(t, err)

	existingPayment, err := paymentDB.GetByID(payment.ID)
	assert.NoError(t, err)
	assert.NotNil(t, existingPayment)

	err = paymentDB.Delete(existingPayment.ID)
	assert.NoError(t, err)

	deletedPayment, err := paymentDB.GetByID(payment.ID)
	assert.Error(t, err)
	assert.Nil(t, deletedPayment)
}

func TestGetByOrderID_ValidOrderID_ReturnsPaymentSuccessfully(t *testing.T) {
	db := repository.SetupTestDB(model.Payment{})
	paymentDB := NewPayment(db)

	payment, err := model.NewPayment(1, "stripe", "credit_card", "CAD", "pending", 3189)
	assert.NoError(t, err)
	err = paymentDB.Create(payment)
	assert.NoError(t, err)

	result, err := paymentDB.GetByOrderID(payment.OrderID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	
	// Detailed validation
	assert.Equal(t, payment.OrderID, result.OrderID)
	assert.Equal(t, "CAD", result.Currency)
	assert.Equal(t, "stripe", result.Provider)
	assert.Equal(t, int64(3189), result.AmountCents)
	assert.Equal(t, "pending", result.Status)
	assert.Equal(t, "credit_card", result.Method)
}

func TestGetByOrderID_InvalidOrderID_ReturnsError(t *testing.T) {
	db := repository.SetupTestDB(model.Payment{}, model.Order{})
	paymentDB := NewPayment(db)

	tests := []struct{
		name string
		orderID int
	}{
		{
			name: "Zero Order ID",
			orderID: 0,
		},
		{
			name: "Non-existent Order ID",
			orderID: 9999,
		},
		{
			name: "Negative Order ID",
			orderID: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := paymentDB.GetByOrderID(tt.orderID)

			assert.Error(t, err)
			assert.Nil(t, result)
		})
	}
}