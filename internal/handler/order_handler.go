package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
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

func (h *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 0)
	if err != nil {
		http.Error(w, "Invalid Id", http.StatusBadRequest)
		return
	}
	order, err := h.OrderDB.GetOrderByID(int(id))
	if err != nil {
		http.Error(w, "Canno return order", http.StatusNotFound)
		return
	}
	w.Header().Set("Contenty-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&order)
}

func (h *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 0)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}
	_, err = h.OrderDB.GetOrderByID(int(id))
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}
	var orderReq dto.OrderDTO
	err = json.NewDecoder(r.Body).Decode(&orderReq)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	items := dto.ToOrderItems(orderReq.Items)
	payments := dto.ToPayments(orderReq.Payments)

	order := &model.Order{
		ID: int(id),
		UserID: orderReq.UserID,
		Status: orderReq.Status,
		TotalPrice: orderReq.TotalPrice,
		Currency: orderReq.Currency,
		Items: items,
		Payments: payments,
	}

	err = h.OrderDB.UpdateOrder(order)
	if err != nil {
		http.Error(w, "Error updating order", http.StatusInternalServerError)
		return
	}

	updatedOrder, err := h.OrderDB.GetOrderByID(order.ID)
	if err != nil {
		http.Error(w, "Cannot return order", http.StatusBadRequest)
		return
	}

	orderDTO := dto.ToOrderDTO(updatedOrder)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orderDTO)
}

func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 0)
	if err != nil {
		http.Error(w, "Invalid Id", http.StatusBadRequest)
		return
	}

	_, err = h.OrderDB.GetOrderByID(int(id))
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	err = h.OrderDB.DeleteOrder(int(id))
	if err != nil {
		http.Error(w, "Order cannot be delete", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}