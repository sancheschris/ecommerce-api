package service

import (
	"errors"
	"testing"

	"github.com/sancheschris/ecommerce-api/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	stripe "github.com/stripe/stripe-go/v76"
)

func TestCreatePayment_Success(t *testing.T) {
	// Arrange
	mockRepo := &MockPaymentRepo{}
	mockClient := &MockPaymentClient{}
	service := NewPaymentService(mockRepo, mockClient)

	orderID := 1
	amount := int64(2000)
	currency := "usd"
	intentID := "pi_123456"

	// mock stripe response
	intent := &stripe.PaymentIntent{
		ID:       intentID,
		Status:   stripe.PaymentIntentStatusRequiresPaymentMethod,
		Amount:   amount,
		Currency: stripe.Currency(currency),
	}

	mockClient.On("CreatePaymentIntent", amount, currency).Return(intent, nil)
	mockRepo.On("Create", mock.AnythingOfType("*model.Payment")).Return(nil)

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
	mockClient.AssertExpectations(t)
	mockRepo.AssertNumberOfCalls(t, "Create", 1)
	mockClient.AssertNumberOfCalls(t, "CreatePaymentIntent", 1)
}

func TestCreatePayment__StripeError(t *testing.T) {
	// Arrange
	mockRepo := &MockPaymentRepo{}
	mockClient := &MockPaymentClient{}
	service := NewPaymentService(mockRepo, mockClient)

	orderID := 1
	amount := int64(2000)
	currency := "usd"

	stripeError := errors.New("stripe connection failed")
	mockClient.On("CreatePaymentIntent", amount, currency).Return((*stripe.PaymentIntent)(nil), stripeError)

	// Act
	result, err := service.CreatePayment(orderID, amount, currency)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to create stripe payment intent")
	assert.Contains(t, err.Error(), "stripe connection failed")

	mockClient.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Create")
}

func TestUpdatePaymentStatus_NilStripePaymentIntentID(t *testing.T) {
	// Arrange
	repo := new(MockPaymentRepo)
	client := new(MockPaymentClient)
	service := NewPaymentService(repo, client)

	paymentID := 1

	// Create a payment with nil StripePaymentIntentID
	paymentWithoutStripeID := &model.Payment{
		ID:                    paymentID,
		OrderID:               123,
		AmountCents:           2000,
		Currency:              "usd",
		Status:                "pending",
		StripePaymentIntentID: nil, // This is the key - nil pointer
		Method:                "stripe",
	}

	repo.On("GetByID", paymentID).Return(paymentWithoutStripeID, nil)

	// Act
	result, err := service.UpdatePaymentStatus(paymentID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "payment has no stripe payment intent ID")

	// Verify that we never called Stripe since we returned early
	client.AssertNotCalled(t, "GetPaymentIntentStatus")
	repo.AssertExpectations(t)
}

func TestCreatePayment_DatabaseError(t *testing.T) {
	// Arrange
	mockRepo := &MockPaymentRepo{}
	mockClient := &MockPaymentClient{}
	service := NewPaymentService(mockRepo, mockClient)

	orderID := 1
	amount := int64(2000)
	currency := "usd"
	intentID := "pi_test_123456789"

	mockIntent := &stripe.PaymentIntent{
		ID:       intentID,
		Amount:   amount,
		Currency: stripe.Currency(currency),
		Status:   stripe.PaymentIntentStatusRequiresPaymentMethod,
	}

	dbError := errors.New("database connection failed")
	mockClient.On("CreatePaymentIntent", amount, currency).Return(mockIntent, nil)
	mockRepo.On("Create", mock.AnythingOfType("*model.Payment")).Return(dbError)

	// Act
	result, err := service.CreatePayment(orderID, amount, currency)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to save payment to database")
	assert.Contains(t, err.Error(), "database connection failed")

	mockClient.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestCreatePayment_DifferentCurrencies(t *testing.T) {
	testCases := []struct {
		name     string
		currency string
	}{
		{"USD", "usd"},
		{"EUR", "eur"},
		{"GBP", "gbp"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockRepo := &MockPaymentRepo{}
			mockClient := &MockPaymentClient{}
			service := NewPaymentService(mockRepo, mockClient)

			orderID := 1
			amount := int64(1500)
			intentID := "pi_test_" + tc.currency

			mockIntent := &stripe.PaymentIntent{
				ID:       intentID,
				Amount:   amount,
				Currency: stripe.Currency(tc.currency),
				Status:   stripe.PaymentIntentStatusRequiresPaymentMethod,
			}

			mockClient.On("CreatePaymentIntent", amount, tc.currency).Return(mockIntent, nil)
			mockRepo.On("Create", mock.AnythingOfType("*model.Payment")).Return(nil)

			// Act
			result, err := service.CreatePayment(orderID, amount, tc.currency)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tc.currency, result.Currency)

			mockClient.AssertExpectations(t)
			mockRepo.AssertExpectations(t)
		})
	}
}
