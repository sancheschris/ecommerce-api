package product

import (
	"github.com/sancheschris/ecommerce-api/internal/model"
	"gorm.io/gorm"
)

type Product struct {
	DB *gorm.DB
}

func NewProduct(db *gorm.DB) *Product {
	return &Product{DB: db}
}

func (p *Product) Create(product *model.Product) error {
	return p.DB.Create(product).Error
}