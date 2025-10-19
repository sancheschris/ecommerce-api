package stripe

import (
	"errors"
	"fmt"

	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/customer"
	"github.com/stripe/stripe-go/v79/paymentintent"
	"github.com/stripe/stripe-go/v79/paymentmethod"
)

type Client struct {
	secretKey string
}

type PaymentIntentRequest struct {
	Amount        int64  `json:"amount"`
	Currency      string `json:"currency"`
	CustomerID    string `json:"customer_id,omitempty"`
	PaymentMethod string `json:"payment_method,omitempty"`
	Description   string `json:"description,omitempty"`
}

type PaymentIntentResponse struct {
	ID           string `json:"id"`
	ClientSecret string `json:"client_secret"`
	Status       string `json:"status"`
	Amount       int64  `json:"amount"`
	Currency     string `json:"currency"`
}

type CustomerRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CustomerResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func NewClient(secretKey string) *Client {
	stripe.Key = secretKey
	return &Client{
		secretKey: secretKey,
	}
}

// CreateCustomer creates a new Stripe customer
func (c *Client) CreateCustomer(req CustomerRequest) (*CustomerResponse, error) {
	params := &stripe.CustomerParams{
		Name:  stripe.String(req.Name),
		Email: stripe.String(req.Email),
	}

	customer, err := customer.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

	return &CustomerResponse{
		ID:    customer.ID,
		Name:  customer.Name,
		Email: customer.Email,
	}, nil
}

// CreatePaymentIntent creates a new payment intent
func (c *Client) CreatePaymentIntent(req PaymentIntentRequest) (*PaymentIntentResponse, error) {
	if req.Amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}

	if req.Currency == "" {
		req.Currency = "usd" // default currency
	}

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(req.Amount),
		Currency: stripe.String(req.Currency),
	}

	if req.CustomerID != "" {
		params.Customer = stripe.String(req.CustomerID)
	}

	if req.PaymentMethod != "" {
		params.PaymentMethod = stripe.String(req.PaymentMethod)
		params.ConfirmationMethod = stripe.String("manual")
		params.Confirm = stripe.Bool(true)
	}

	if req.Description != "" {
		params.Description = stripe.String(req.Description)
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment intent: %w", err)
	}

	return &PaymentIntentResponse{
		ID:           pi.ID,
		ClientSecret: pi.ClientSecret,
		Status:       string(pi.Status),
		Amount:       pi.Amount,
		Currency:     string(pi.Currency),
	}, nil
}

// ConfirmPaymentIntent confirms a payment intent
func (c *Client) ConfirmPaymentIntent(paymentIntentID string, paymentMethodID string) (*PaymentIntentResponse, error) {
	params := &stripe.PaymentIntentConfirmParams{
		PaymentMethod: stripe.String(paymentMethodID),
	}

	pi, err := paymentintent.Confirm(paymentIntentID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm payment intent: %w", err)
	}

	return &PaymentIntentResponse{
		ID:           pi.ID,
		ClientSecret: pi.ClientSecret,
		Status:       string(pi.Status),
		Amount:       pi.Amount,
		Currency:     string(pi.Currency),
	}, nil
}

// GetPaymentIntent retrieves a payment intent
func (c *Client) GetPaymentIntent(paymentIntentID string) (*PaymentIntentResponse, error) {
	pi, err := paymentintent.Get(paymentIntentID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment intent: %w", err)
	}

	return &PaymentIntentResponse{
		ID:           pi.ID,
		ClientSecret: pi.ClientSecret,
		Status:       string(pi.Status),
		Amount:       pi.Amount,
		Currency:     string(pi.Currency),
	}, nil
}

// AttachPaymentMethodToCustomer attaches a payment method to a customer
func (c *Client) AttachPaymentMethodToCustomer(paymentMethodID, customerID string) error {
	params := &stripe.PaymentMethodAttachParams{
		Customer: stripe.String(customerID),
	}

	_, err := paymentmethod.Attach(paymentMethodID, params)
	if err != nil {
		return fmt.Errorf("failed to attach payment method to customer: %w", err)
	}

	return nil
}