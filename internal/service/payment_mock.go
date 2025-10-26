package service

import (
	"github.com/stretchr/testify/mock"
	stripe "github.com/stripe/stripe-go/v76"
)

type MockPaymentClient struct {
	mock.Mock
}

func (m *MockPaymentClient) CreatePaymentIntent(amount int64, currency string) (*stripe.PaymentIntent, error) {
	args := m.Called(amount, currency)
	return args.Get(0).(*stripe.PaymentIntent), args.Error(1)
}

func (m *MockPaymentClient) ConfirmPaymentIntent(paymentIntent string) (*stripe.PaymentIntent, error) {
	args := m.Called(paymentIntent)
	return args.Get(0).(*stripe.PaymentIntent), args.Error(1)
}

func (m *MockPaymentClient) GetPaymentIntentStatus(PaymentIntentID string) (*stripe.PaymentIntent, error) {
	args := m.Called(PaymentIntentID)
	return args.Get(0).(*stripe.PaymentIntent), args.Error(1)
}

func (m *MockPaymentClient) CreateCustomer(email, name string) (*stripe.Customer, error) {
	args := m.Called(email, name)
	return args.Get(0).(*stripe.Customer),args.Error(1)
}

type MockPaymentRepo struct {
	mock.Mock
}



