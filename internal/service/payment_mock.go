package service

import (
	"github.com/sancheschris/ecommerce-api/internal/model"
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

func (m *MockPaymentClient) ConfirmPaymentIntent(paymentIntentID string) (*stripe.PaymentIntent, error) {
    args := m.Called(paymentIntentID)
    return args.Get(0).(*stripe.PaymentIntent), args.Error(1)
}

func (m *MockPaymentClient) GetPaymentIntentStatus(paymentIntentID string) (*stripe.PaymentIntent, error) {
    args := m.Called(paymentIntentID)
    return args.Get(0).(*stripe.PaymentIntent), args.Error(1)
}

func (m *MockPaymentClient) CreateCustomer(email, name string) (*stripe.Customer, error) {
    args := m.Called(email, name)
    return args.Get(0).(*stripe.Customer), args.Error(1) // Fixed spacing
}

type MockPaymentRepo struct {
    mock.Mock
}

func (m *MockPaymentRepo) Create(payment *model.Payment) error {
    args := m.Called(payment)
    return args.Error(0)
}

func (m *MockPaymentRepo) GetByID(id int) (*model.Payment, error) {
    args := m.Called(id)
    return args.Get(0).(*model.Payment), args.Error(1)
}

func (m *MockPaymentRepo) Update(payment *model.Payment) error {
    args := m.Called(payment)
    return args.Error(0)
}

func (m *MockPaymentRepo) Delete(id int) error {
    args := m.Called(id)
    return args.Error(0)
}

func (m *MockPaymentRepo) GetByOrderID(orderID int) (*model.Payment, error) {
    args := m.Called(orderID)
    return args.Get(0).(*model.Payment), args.Error(1)
}

func (m *MockPaymentRepo) GetByUserID(userID int) (*model.Payment, error) {
    args := m.Called(userID)
    return args.Get(0).(*model.Payment), args.Error(1)
}

func (m *MockPaymentRepo) GetByStatus(status string) ([]*model.Payment, error) {
    args := m.Called(status)
    return args.Get(0).([]*model.Payment), args.Error(1)
}