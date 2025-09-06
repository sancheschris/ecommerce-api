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

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var orderRequest dto.OrderDTO
	err := json.NewDecoder(r.Body).Decode(&orderRequest)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return	
	}

	items := dto.ToOrderItems(orderRequest.Items)
	payments := dto.ToPayments(orderRequest.Payments)

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

	savedOrder, err := h.OrderDB.GetOrderByID(o.ID)
	if err != nil {
        http.Error(w, "Error fetching created order", http.StatusInternalServerError)
        return
    }

	orderDTO := dto.ToOrderDTO(savedOrder)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(orderDTO)
}

func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.OrderDB.GetOrders()
	if err != nil {
		http.Error(w, "Error to return orders", http.StatusInternalServerError)
		return
	}
	ordersDTO := make([]dto.OrderDTO, len(orders))
	for i, o := range orders {
		ordersDTO[i] = dto.ToOrderDTO(&o)
	}
	json.NewEncoder(w).Encode(&orders)
	w.WriteHeader(http.StatusOK)
}