package handler

import (
	"encoding/json"
	"net/http"

	"github.com/sancheschris/ecommerce-api/internal/dto"
	repo "github.com/sancheschris/ecommerce-api/internal/repository/order"
)

type OrderHandler struct {
	OrderDB repo.OrderInterface 
}

func NewHandler(orderDB repo.OrderInterface) *OrderHandler {
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
}