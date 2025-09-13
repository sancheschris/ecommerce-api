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
		name string
		orderID int
		provider string
		amountCents int64
		method string
		currency string
		status string
		wantErr bool
	}{
		{
			name: "Valid Payment",
			orderID: 1,
			provider: "stripe",
			amountCents: 10,
			method: "credit_card",
			currency: "USD",
			status: "pending",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payment, err := model.NewPayment(tt.orderID, tt.provider, tt.method, tt.currency, tt.status, tt.amountCents)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, payment)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, payment)
				err = paymentDB.Create(payment)
				assert.NoError(t, err)
			}
		})
	}
}

