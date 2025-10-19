package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/sancheschris/ecommerce-api/internal/model"
	paymentRepo "github.com/sancheschris/ecommerce-api/internal/repository/payment"
	stripe_client "github.com/sancheschris/ecommerce-api/pkg/stripe"
)

type PaymentHandler struct {
	PaymentDB    paymentRepo.PaymentInterface
	StripeClient *stripe_client.Client
}

type CreatePaymentRequest struct {
	OrderID       int    `json:"order_id"`
	PaymentMethod string `json:"payment_method"`
	Amount        int64  `json:"amount"`
	Currency      string `json:"currency"`
	CustomerID    string `json:"customer_id,omitempty"`
}

type PaymentResponse struct {
	ID                    int    `json:"id"`
	OrderID               int    `json:"order_id"`
	Provider              string `json:"provider"`
	AmountCents           int64  `json:"amount_cents"`
	Method                string `json:"method"`
	Currency              string `json:"currency"`
	Status                string `json:"status"`
	StripePaymentIntentID string `json:"stripe_payment_intent_id,omitempty"`
	ClientSecret          string `json:"client_secret,omitempty"`
}

type CreateCustomerRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func NewPaymentHandler(paymentDB paymentRepo.PaymentInterface, stripeClient *stripe_client.Client) *PaymentHandler {
	return &PaymentHandler{
		PaymentDB:    paymentDB,
		StripeClient: stripeClient,
	}
}

// CreatePayment godoc
// @Summary Create a payment
// @Description Create a new payment and process it with Stripe
// @Tags payments
// @Accept json
// @Produce json
// @Param request body CreatePaymentRequest true "Payment request"
// @Success 201 {object} PaymentResponse
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /payments [post]
// @Security ApiKeyAuth
func (h *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	var req CreatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create Stripe Payment Intent
	stripeReq := stripe_client.PaymentIntentRequest{
		Amount:      req.Amount,
		Currency:    req.Currency,
		Description: "Order payment",
	}

	// Only set customer ID if provided and not empty
	if req.CustomerID != "" {
		stripeReq.CustomerID = req.CustomerID
	}

	stripeResp, err := h.StripeClient.CreatePaymentIntent(stripeReq)
	if err != nil {
		http.Error(w, "Failed to create payment intent: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Create payment record in database
	payment, err := model.NewPayment(
		req.OrderID,
		"stripe",
		req.PaymentMethod,
		req.Currency,
		"pending",
		req.Amount,
	)
	if err != nil {
		http.Error(w, "Invalid payment data: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Set Stripe payment intent ID
	payment.StripePaymentIntentID = &stripeResp.ID

	// Save to database
	if err := h.PaymentDB.Create(payment); err != nil {
		http.Error(w, "Failed to save payment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	response := PaymentResponse{
		ID:                    payment.ID,
		OrderID:               payment.OrderID,
		Provider:              payment.Provider,
		AmountCents:           payment.AmountCents,
		Method:                payment.Method,
		Currency:              payment.Currency,
		Status:                payment.Status,
		StripePaymentIntentID: *payment.StripePaymentIntentID,
		ClientSecret:          stripeResp.ClientSecret,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// ConfirmPayment godoc
// @Summary Confirm a payment
// @Description Confirm a payment with Stripe
// @Tags payments
// @Accept json
// @Produce json
// @Param id path string true "Payment ID"
// @Param request body map[string]string true "Confirmation request with payment_method_id"
// @Success 200 {object} PaymentResponse
// @Failure 400 {object} Error
// @Failure 404 {object} Error
// @Failure 500 {object} Error
// @Router /payments/{id}/confirm [post]
// @Security ApiKeyAuth
func (h *PaymentHandler) ConfirmPayment(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid payment ID", http.StatusBadRequest)
		return
	}

	var req map[string]string
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	paymentMethodID, ok := req["payment_method_id"]
	if !ok {
		http.Error(w, "payment_method_id is required", http.StatusBadRequest)
		return
	}

	// Get payment from database
	payment, err := h.PaymentDB.GetByID(id)
	if err != nil {
		http.Error(w, "Payment not found", http.StatusNotFound)
		return
	}

	if payment.StripePaymentIntentID == nil {
		http.Error(w, "No Stripe payment intent found for this payment", http.StatusBadRequest)
		return
	}

	// Confirm with Stripe
	stripeResp, err := h.StripeClient.ConfirmPaymentIntent(*payment.StripePaymentIntentID, paymentMethodID)
	if err != nil {
		http.Error(w, "Failed to confirm payment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Update payment status
	payment.Status = stripeResp.Status
	if err := h.PaymentDB.Update(payment); err != nil {
		http.Error(w, "Failed to update payment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	response := PaymentResponse{
		ID:                    payment.ID,
		OrderID:               payment.OrderID,
		Provider:              payment.Provider,
		AmountCents:           payment.AmountCents,
		Method:                payment.Method,
		Currency:              payment.Currency,
		Status:                payment.Status,
		StripePaymentIntentID: *payment.StripePaymentIntentID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetPayment godoc
// @Summary Get a payment
// @Description Get payment details by ID
// @Tags payments
// @Accept json
// @Produce json
// @Param id path string true "Payment ID"
// @Success 200 {object} PaymentResponse
// @Failure 400 {object} Error
// @Failure 404 {object} Error
// @Router /payments/{id} [get]
// @Security ApiKeyAuth
func (h *PaymentHandler) GetPayment(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid payment ID", http.StatusBadRequest)
		return
	}

	payment, err := h.PaymentDB.GetByID(id)
	if err != nil {
		http.Error(w, "Payment not found", http.StatusNotFound)
		return
	}

	response := PaymentResponse{
		ID:          payment.ID,
		OrderID:     payment.OrderID,
		Provider:    payment.Provider,
		AmountCents: payment.AmountCents,
		Method:      payment.Method,
		Currency:    payment.Currency,
		Status:      payment.Status,
	}

	if payment.StripePaymentIntentID != nil {
		response.StripePaymentIntentID = *payment.StripePaymentIntentID
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// CreateCustomer godoc
// @Summary Create a Stripe customer
// @Description Create a new customer in Stripe
// @Tags payments
// @Accept json
// @Produce json
// @Param request body CreateCustomerRequest true "Customer request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /customers [post]
// @Security ApiKeyAuth
func (h *PaymentHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var req CreateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create customer in Stripe
	stripeReq := stripe_client.CustomerRequest{
		Name:  req.Name,
		Email: req.Email,
	}

	customer, err := h.StripeClient.CreateCustomer(stripeReq)
	if err != nil {
		http.Error(w, "Failed to create customer: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(customer)
}