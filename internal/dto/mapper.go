package dto

import (
	"time"

	"github.com/sancheschris/ecommerce-api/internal/model"
)

func ToOrderItems(reqItems []OrderItemDTO) []model.OrderItem {
	items := make([]model.OrderItem, len(reqItems))
	for i, item := range reqItems {
		items[i] = model.OrderItem{
			ID:        item.ID,       
			ProductID: item.ProductID,
			Qty:       item.Qty,
			UnitPrice: item.UnitPrice,
		}
	}
	return items
}

func ToPayments(reqPayments []PaymentDTO) []model.Payment {
	payments := make([]model.Payment, len(reqPayments))
	for i, payment := range reqPayments {
		payments[i] = model.Payment{
			ID:          payment.ID,  
			Provider:    payment.Provider,
			AmountCents: int64(payment.Amount),
			Method:      payment.Method,
			Currency:    payment.Currency,
			Status:      payment.Status,
		}
	}
	return payments
}

func ToOrderDTO(order *model.Order) OrderDTO {
	items := make([]OrderItemDTO, len(order.Items))
	for i, item := range order.Items {
		// Map product details if available
		productDTO := ProductDTO{}
		if item.Product != nil && item.Product.ID != 0 {
			productDTO = ProductDTO{
				ID:        int(item.Product.ID),
				Name:      item.Product.Name,
				Price:     item.Product.Price,
				Active:    item.Product.Active,
				CreatedAt: item.Product.CreatedAt.Format(time.RFC3339),
				UpdatedAt: item.Product.UpdatedAt.Format(time.RFC3339),
			}
		}

		items[i] = OrderItemDTO{
			ID:        item.ID,
			OrderID:   item.OrderID,
			ProductID: item.ProductID,
			Product:   productDTO,
			Qty:       item.Qty,
			UnitPrice: item.UnitPrice,
		}
	}
	payments := make([]PaymentDTO, len(order.Payments))
	for i, p := range order.Payments {
		payments[i] = PaymentDTO{
			ID:          p.ID,
			OrderID:     p.OrderID,
			Provider:    p.Provider,
			Amount: float64(p.AmountCents),
			Method:      p.Method,
			Currency:    p.Currency,
			Status:      p.Status,
			CreatedAt:   p.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   p.UpdatedAt.Format(time.RFC3339),
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