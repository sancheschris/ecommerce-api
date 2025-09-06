package dto

import (
	"github.com/sancheschris/ecommerce-api/internal/model"
)

func ToOrderItems(reqItems []OrderItemRequest) []model.OrderItem {
	items := make([]model.OrderItem, len(reqItems))
	for i, item := range reqItems {
		items[i] = model.OrderItem{
			ProductID: item.ProductID,
			Qty: item.Qty,
			UnitPrice: item.UnitPrice,
		}
	}
	return items
}

func ToPayments(reqPayments []PaymentRequest) []model.Payment {
	payments := make([]model.Payment, len(reqPayments))
	for i, payment := range reqPayments {
		payments[i] = model.Payment{
			Provider: payment.Provider,
			AmountCents: payment.Amount,
			Method: payment.Method,
			Currency: payment.Currency,
			Status: payment.Status,
		}
	}
	return payments
}