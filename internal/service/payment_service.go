package service

import (
	"fmt"

	"github.com/sancheschris/ecommerce-api/internal/model"
	paymentRepo "github.com/sancheschris/ecommerce-api/internal/repository/payment"
	"github.com/sancheschris/ecommerce-api/pkg/payment"
)

type PaymentService struct {
	paymentRepo *paymentRepo.Payment
	stripeClient payment.PaymentClient
}

func NewPaymentService(repo *paymentRepo.Payment, client payment.PaymentClient) *PaymentService {
	return &PaymentService{
		paymentRepo: repo,
		stripeClient: client,
	}
}

func (p *PaymentService) CreatePayment(orderID int, amount int64, currency string) (*model.Payment, error) {
	intent, err := p.stripeClient.CreatePaymentIntent(amount, currency)
	if err != nil {
		return nil, fmt.Errorf("failed to create stripe payment intent: %w", err)
	}

	payment := &model.Payment{
		OrderID: orderID,
		AmountCents: amount / 100,
		Currency: currency,
		Status: string(intent.Status),
		StripePaymentIntentID: &intent.ID,
		Method: "stripe",
	}

	if err := p.paymentRepo.Create(payment); err != nil {
		return nil, fmt.Errorf("failed to save payment to database: %w", err)
	}

	return payment, nil
}

func (p *PaymentService) UpdatePaymentStatus(paymentID int) (*model.Payment, error) {
	payment, err := p.paymentRepo.GetByID(paymentID)
	if err != nil {
		return nil, fmt.Errorf("payment not found: %w", err)
	}

	// check status with stripe
	intent, err := p.stripeClient.GetPaymentIntentStatus(*payment.StripePaymentIntentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get stripe status: %w", err)
	}

	// update database if status changed
	 newStatus := string(intent.Status)
    if payment.Status != newStatus {
        payment.Status = newStatus
        if err := p.paymentRepo.Update(payment); err != nil {
            return nil, fmt.Errorf("failed to update payment status: %w", err)
        }
    }
    return payment, nil
}

func (p *PaymentService) GetPaymentByOrderID(orderID int) (*model.Payment, error) {
	return p.paymentRepo.GetByOrderID(orderID)
}

func (p *PaymentService) GetPaymentsByStatus(status string) ([]*model.Payment, error) {
	return p.paymentRepo.GetByStatus(status)
}
