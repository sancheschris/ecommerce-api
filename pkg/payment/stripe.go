package payment

import (
	stripe "github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/client"
)

type PaymentClient interface {
	CreatePaymentIntent(amount int64, currency string) (*stripe.PaymentIntent, error)
	ConfirmPaymentIntent(paymentIntentID string) (*stripe.PaymentIntent, error)
	GetPaymentIntentStatus(paymentIntentID string) (string, error)
	CreateCustomer(email, name string) (*stripe.Customer, error)
}

type StripeClient struct {
	client *client.API
	
}

func NewStripeClient(secretKey string) *StripeClient {
	sc := &client.API{}
	sc.Init(secretKey, nil)
	return &StripeClient{client: sc}
}

func (s *StripeClient) CreatePaymentIntent(amount int64, currency string) (*stripe.PaymentIntent, error) {
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amount),
		Currency: stripe.String(currency),
	}
	return s.client.PaymentIntents.New(params)
}

func (s *StripeClient) ConfirmPaymentIntent(paymentIntentID string) (*stripe.PaymentIntent, error) {
	params := &stripe.PaymentIntentConfirmParams{}
	return s.client.PaymentIntents.Confirm(paymentIntentID, params)
}

func (s *StripeClient) GetPaymentIntentStatus(paymentIntentID string) (*stripe.PaymentIntent, error) {
	return s.client.PaymentIntents.Get(paymentIntentID, nil)
}

func (s *StripeClient) CreateCustomer(email, name string) (*stripe.Customer, error) {
	params := &stripe.CustomerParams{
		Email: stripe.String(email),
		Name:  stripe.String(name),
	}
	return s.client.Customers.New(params)
}

