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

func (p *Product) GetProductByID(id int64) (*model.Product, error) {
	var product model.Product
	err := p.DB.First(&product, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &product, err
}

func (p *Product) GetProducts() ([]model.Product, error) {
	var products []model.Product
	err := p.DB.Find(&products).Error
	return products, err
}

