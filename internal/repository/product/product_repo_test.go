package product

import (
	"testing"

	"github.com/sancheschris/ecommerce-api/internal/model"
	"github.com/sancheschris/ecommerce-api/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestNewProduct(t *testing.T) {
	db := repository.SetupTestDB(model.Product{})
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

func TestGetProductByID(t *testing.T) {
	db := repository.SetupTestDB(model.Product{})
	productDB := NewProduct(db)

	product, err := model.NewProduct("Macbook M4 Pro", 26000.00, true)
	assert.NoError(t, err)
	assert.NotNil(t, product)

	db.Create(product)

	product, err = productDB.GetProductByID(product.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Macbook M4 Pro", product.Name)
}

func TestGetProductByIdTable(t *testing.T) {
	db := repository.SetupTestDB(model.Product{})
	productDB := NewProduct(db)

	product1, _ := model.NewProduct("Macbook M4 Pro", 26000.00, true)
    db.Create(product1)
    product2, _ := model.NewProduct("Iphone 16", 8900.00, true)
    db.Create(product2)

	tests := []struct {
    name     string
    id       int64
    wantName string
    wantPrice float64
    wantErr  bool
}{
    {"Valid ID", product2.ID, "Iphone 16", 8900.00, false},
    {"Invalid ID", -1, "", 0.0, true},
    {"Non-existent ID", 9999, "", 0.0, true},
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        product, err := productDB.GetProductByID(tt.id)
        if tt.wantErr {
            assert.Error(t, err)
            assert.Nil(t, product)
        } else {
            assert.NoError(t, err)
            assert.NotNil(t, product)
            assert.Equal(t, tt.wantName, product.Name)
            assert.Equal(t, tt.wantPrice, product.Price)
        }
    })
	}
}

func TestGetProducts(t *testing.T) {
	db := repository.SetupTestDB(model.Product{})
	productDB := NewProduct(db)

	p1, _ := model.NewProduct("Tablet", 3500.00, true)
	p2, _ := model.NewProduct("TV", 5400.00, true)

	db.Create(p1)
	db.Create(p2)

	products, err := productDB.GetProducts()
	assert.NoError(t, err)
	assert.Len(t, products, 2)

	names := []string{products[0].Name, products[1].Name}

	assert.Contains(t, names, "Tablet")
	assert.Contains(t, names, "TV")
}

func TestDeleteProducts(t *testing.T) {
	db := repository.SetupTestDB(model.Product{})
	productDB := NewProduct(db)

	p1, _ := model.NewProduct("Tablet", 3500.00, true)
	db.Create(p1)

	err := productDB.Delete(p1.ID)
	assert.NoError(t, err)

	product, err := productDB.GetProductByID(p1.ID)
	assert.Error(t, err)
	assert.Nil(t, product)
}

func TestUpdateProduct(t *testing.T) {
	db := repository.SetupTestDB(model.Product{})
	productDB := NewProduct(db)

	p1, _ := model.NewProduct("Iphone 15", 3600.00, true)
	db.Create(p1)
	
	currentProduct, err := productDB.GetProductByID(p1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, currentProduct)
	assert.Equal(t, "Iphone 15", currentProduct.Name)

	currentProduct.Price = 5800.00
	currentProduct.Name = "Updated Product"

	err = productDB.Update(currentProduct)
	assert.NoError(t, err)

	updatedProduct, err := productDB.GetProductByID(currentProduct.ID)
	assert.NoError(t, err)
	assert.Equal(t, 5800.00, updatedProduct.Price)
	assert.Equal(t, "Updated Product", updatedProduct.Name)

}