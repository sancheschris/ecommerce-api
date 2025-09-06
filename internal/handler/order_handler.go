package handler

import (
	"encoding/json"
	"net/http"

	"github.com/sancheschris/ecommerce-api/internal/dto"
	"github.com/sancheschris/ecommerce-api/internal/model"
	repo "github.com/sancheschris/ecommerce-api/internal/repository/order"
)

type OrderHandler struct {
	OrderDB repo.OrderInterface 
}

func NewOrderHandler(orderDB repo.OrderInterface) *OrderHandler {
	return &OrderHandler{
		OrderDB: orderDB,
	}
}

func (h OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var orderRequest dto.OrderRequest
	err := json.NewDecoder(r.Body).Decode(&orderRequest)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return	
	}
	items := make([]model.OrderItem, len(orderRequest.Items))
	for i, item := range orderRequest.Items {
		items[i] = model.OrderItem{
			ProductID: item.ProductID,
			Qty: item.Qty,
			UnitPrice: item.UnitPrice,
		}
	}

	payments := make([]model.Payment, len(orderRequest.Payments))
	for i, payment := range orderRequest.Payments {
		payments[i] = model.Payment{
			Provider: payment.Provider,
			AmountCents: payment.Amount,
			Method: payment.Method,
			Currency: payment.Currency,
			Status: payment.Status,
		}
	}

	o, err := model.NewOrder(orderRequest.UserID, items, orderRequest.Status, orderRequest.TotalPrice, orderRequest.Currency, payments)
	if err != nil {
		http.Error(w, "Error creating new order", http.StatusBadRequest)
		return
	}
	err = h.OrderDB.CreateOrder(o)
	if err != nil {
		http.Error(w, "Error creating order", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&orderRequest)
}