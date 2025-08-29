package product

import "github.com/sancheschris/ecommerce-api/internal/model"

type ProductInterface interface {
	Create(product *model.Product) error
	Update(product *model.Product) error
	Delete(id int64) error
	GetProductByID(id int64) (*model.Product, error)
	GetProducts() ([]model.Product, error)
}