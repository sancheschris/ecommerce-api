package model

import "time"

type Product struct {
	ID int64 `json:"id" gorm:"primaryKey;autoIncrement"`
	Name string `json:"name"`
	Price float64 `json:"price"`
	Active bool `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewProduct(name string, price float64, active bool) *Product {
	return &Product{
		Name: name,
		Price: price,
		Active: active,
		CreatedAt: time.Now(),
	}
}