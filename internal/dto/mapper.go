package dto

import (
	"time"

	"github.com/sancheschris/ecommerce-api/internal/model"
)

func ToOrderItems(reqItems []OrderItemDTO) []model.OrderItem {
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

func ToPayments(reqPayments []PaymentDTO) []model.Payment {
	payments := make([]model.Payment, len(reqPayments))
	for i, payment := range reqPayments {
		payments[i] = model.Payment{
			Provider: payment.Provider,
			AmountCents: payment.AmountCents,
			Method: payment.Method,
			Currency: payment.Currency,
			Status: payment.Status,
		}
	}
	return payments
}

func ToOrderDTO(order *model.Order) OrderDTO {
	items := make([]OrderItemDTO, len(order.Items))
	for i, item := range order.Items {
		items[i] = OrderItemDTO{
			ID: item.ID,
			ProductID: item.ProductID,
			Qty: item.Qty,
			UnitPrice: item.UnitPrice,
		}
	}
	payments := make([]PaymentDTO, len(order.Payments))
	for i, p := range order.Payments {
		payments[i] = PaymentDTO{
			ID: p.ID,
			Provider: p.Provider,
			AmountCents: p.AmountCents,
			Method: p.Method,
			Currency: p.Currency,
			Status: p.Status,
		}
	}
	return OrderDTO{
		ID: order.ID,
		UserID: order.UserID,
		Status: order.Status,
		TotalPrice: order.TotalPrice,
		Currency: order.Currency,
		Items: items,
		Payments: payments,
		 CreatedAt:  order.CreatedAt.Format(time.RFC3339),
        UpdatedAt:  order.UpdatedAt.Format(time.RFC3339),
	}
}