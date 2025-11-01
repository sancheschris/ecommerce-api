package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/sancheschris/ecommerce-api/internal/dto"
	"github.com/sancheschris/ecommerce-api/internal/service"
)

type PaymentHandler struct {
	paymentService *service.PaymentService
}

func NewPaymentHandler(paymentService *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

// Create payment godoc
// @Summary Create Payment
// @Description Create Payments
// @Tags payments
// @Accept json
// @Product json
// @Param request body dto.CreatePaymentRequest true "payment request"
// @Success 201
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /payments [post]
// @Security ApiKeyAuth
func (h *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	var req dto.CreatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	payment, err := h.paymentService.CreatePayment(req.OrderID, req.Amount, req.Currency)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(payment); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetPaymentStatus godoc
// @Summary Get payment by status
// @Description get payment by status
// @Tags payments
// @Accept json
// @Produce json
// @Param id path string true "payment ID" format(string)
// @Success 200 {object} model.Order
// @Failure 400 {object} Error
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /payments/{id}/status [get]
// @Security ApiKeyAuth
func (h *PaymentHandler) GetPaymentStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		http.Error(w, "Payment ID is required", http.StatusBadRequest)
		return
	}

	paymentID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid payment ID", http.StatusBadRequest)
		return
	}

	payment, err := h.paymentService.UpdatePaymentStatus(paymentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(payment); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
