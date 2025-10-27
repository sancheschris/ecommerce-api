package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	stripe "github.com/stripe/stripe-go/v76"
)



func TestCreatePayment_Success(t *testing.T) {
    repo := &MockPaymentRepo{}
    client := &MockPaymentClient{}
    service := NewPaymentService(repo, client)

    orderID := 1
    amount := int64(2000)
    currency := "usd"
    intentID := "pi_123456"

    // mock stripe response
    intent := &stripe.PaymentIntent{
        ID: intentID,
        Status: stripe.PaymentIntentStatusRequiresPaymentMethod,
        Amount: amount,
        Currency: stripe.Currency(currency),
    }

    client.On("CreatePaymentIntent", amount, currency).Return(intent, nil)
    repo.On("Create", mock.AnythingOfType("*model.Payment")).Return(nil)

    // act
    actual, err := service.CreatePayment(orderID, amount, currency)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, actual)
    assert.Equal(t, orderID, actual.OrderID)
    assert.Equal(t, currency, actual.Currency)
    assert.Equal(t, amount/100, actual.AmountCents)
    assert.Equal(t, intentID, *actual.StripePaymentIntentID)
    assert.Equal(t, string(stripe.PaymentIntentStatusRequiresPaymentMethod), actual.Status)
    assert.Equal(t, "stripe", actual.Method)
	
    // Verify
    client.AssertExpectations(t)
    repo.AssertNumberOfCalls(t, "Create", 1)
    client.AssertNumberOfCalls(t, "CreatePaymentIntent", 1)
}