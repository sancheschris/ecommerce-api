package product

import (
	"testing"

	"github.com/sancheschris/ecommerce-api/internal/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestNewProduct(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&model.Product{})
	productDB := NewProduct(db)

	tests := []struct {
		name string
		price float64
		active bool
		want bool
	} {
		{"Iphone 16", 6900.00, true, false},
		{"", 29000.00, true, true},
		{"Macbook M1", 0.0, true, true},
		{"", -10, true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product, err := model.NewProduct(tt.name, tt.price, tt.active)
			if tt.want {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				err = productDB.Create(product)
				assert.NoError(t, err)
			}
		})
	}

}