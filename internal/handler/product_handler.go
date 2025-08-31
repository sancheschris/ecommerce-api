package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/sancheschris/ecommerce-api/internal/dto"
	"github.com/sancheschris/ecommerce-api/internal/model"
	repo "github.com/sancheschris/ecommerce-api/internal/repository/product"
)

type ProductHandler struct {
	ProductDB repo.ProductInterface
}

func NewProductHandler(productDB repo.ProductInterface) *ProductHandler {
	return &ProductHandler{
		ProductDB: productDB,
	}
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var product dto.ProductRequest
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	p, err := model.NewProduct(product.Name, product.Price, product.Active)
	if err != nil {
		http.Error(w, "Error creating new product", http.StatusBadRequest)
		return
	}
	err = h.ProductDB.Create(p)
	if err != nil {
		http.Error(w, "Error creating product", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&product)
}

func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid Id", http.StatusBadRequest)
		return
	}
	product, err := h.ProductDB.GetProductByID(id)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.ProductDB.GetProducts()
	if err != nil {
		http.Error(w, "Cannot retrieve products", http.StatusInternalServerError)
		return 
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
	}

	_, err = h.ProductDB.GetProductByID(id)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}
	
	var productReq dto.ProductRequest
	err = json.NewDecoder(r.Body).Decode(&productReq)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	product := &model.Product{
		ID:     id,
		Name:   productReq.Name,
		Price:  productReq.Price,
		Active: productReq.Active,
	}

	err = h.ProductDB.Update(product)
	if err != nil {
		http.Error(w, "Error updating product", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "ID is invalid", http.StatusBadRequest)
		return
	}
	_, err = h.ProductDB.GetProductByID(id)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}
	err = h.ProductDB.Delete(id)
	if err != nil {
		http.Error(w, "Failed to delete", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}