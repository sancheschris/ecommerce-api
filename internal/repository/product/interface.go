package product

import "github.com/sancheschris/ecommerce-api/internal/model"

type ProductInterface interface {
	Create(*model.Product) error
	Update(*model.Product) error
	GetProductsByID() (*model.Product, error)
	GetProducts() ([]model.Product, error)
}