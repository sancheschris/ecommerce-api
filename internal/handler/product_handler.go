package handler

import (
	"encoding/json"
	"net/http"

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