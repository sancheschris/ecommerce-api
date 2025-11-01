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

// Create order godoc
// @Summary Create order
// @Description Create orders
// @Tags orders
// @Accept json
// @Produce json
// @Param request body dto.CreateOrderRequest true "order request"
// @Success 201
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /orders [post]
// @Security ApiKeyAuth
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var orderRequest dto.CreateOrderRequest
	err := json.NewDecoder(r.Body).Decode(&orderRequest)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	items := dto.ToOrderItemsFromRequest(orderRequest.Items)

	totalPrice := 0.0
	for _, item := range items {
		totalPrice += item.UnitPrice * float64(item.Qty)
	}

	o, err := model.NewOrder(
		orderRequest.UserID,
		items,
		"pending",
		totalPrice,
		"USD",
		[]model.Payment{},
	)

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
	if err := json.NewEncoder(w).Encode(orderDTO); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// ListOrders godoc
// @Summary List orders
// @Description get all orders
// @Tags orders
// @Accept json
// @Produce json
// @Success 200 {array} model.Order
// @Failure 404 {object} Error
// @Failure 500 {object} Error
// @Router /orders [get]
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&ordersDTO); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetOrder godoc
// @Summary Get order
// @Description get order by id
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "order ID" format(string)
// @Success 200 {object} model.Order
// @Failure 400 {object} Error
// @Failure 404 {object} Error
// @Failure 500 {object} Error
// @Router /orders/{id} [get]
// @Security ApiKeyAuth
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(&order)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// UpdateOrder godoc
// @Summary Update order
// @Description Update order
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "order ID" format(string)
// @Param request body dto.UpdateOrderRequest true "order request"
// @Success 200 {object} dto.OrderDTO
// @Failure 404 {object} Error
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /orders/{id} [put]
// @Security ApiKeyAuth
func (h *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 0)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}
	
	existingOrder, err := h.OrderDB.GetOrderByID(int(id))
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	var updateReq dto.UpdateOrderRequest
	err = json.NewDecoder(r.Body).Decode(&updateReq)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if updateReq.Status != "" {
		existingOrder.Status = updateReq.Status
	}

	if (len(updateReq.Items) > 0) {
		items := dto.ToOrderItemsFromUpdateRequest(updateReq.Items)
		existingOrder.Items = items

		totalPrice := 0.0
		for _, item := range items {
			totalPrice += item.UnitPrice *float64(item.Qty)
		}
		existingOrder.TotalPrice = totalPrice
	}

	err = h.OrderDB.UpdateOrder(existingOrder)
	if err != nil {
		http.Error(w, "Error updating order", http.StatusInternalServerError)
		return
	}

	updatedOrder, err := h.OrderDB.GetOrderByID(existingOrder.ID)
	if err != nil {
		http.Error(w, "Cannot return order", http.StatusBadRequest)
		return
	}

	orderDTO := dto.ToOrderDTO(updatedOrder)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(orderDTO); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// DeleteOrder godoc
// @Summary Delete order
// @Description delete order by id
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "order ID" format(string)
// @Success 204
// @Faikure 400 {object} Error
// @Failure 404 {object} Error
// @Failure 500 {object} Error
// @Router /orders [delete]
// @Security ApiKeyAuth
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

// GetOrderByUserID godoc
// @Summary Get orders by user id
// @Description get orders by user id
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "order ID" format(string)
// @Success 200 {object} dto.OrderDTO
// @Faikure 400 {object} Error
// @Failure 404 {object} Error
// @Failure 500 {object} Error
// @Router /orders [get]
// @Security ApiKeyAuth
func (h *OrderHandler) GetOrdersByUserID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 0)
	if err != nil {
		http.Error(w, "Invalid Id", http.StatusBadRequest)
		return
	}
	ordersByUser, err := h.OrderDB.GetOrdersByUserID(int(id))
	if err != nil {
		http.Error(w, "Cannot get orders by user id", http.StatusNotFound)
		return
	}

	ordersDTO := make([]dto.OrderDTO, len(ordersByUser))
	for i, o := range ordersByUser {
		ordersDTO[i] = dto.ToOrderDTO(&o)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(ordersDTO); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
