package model

import (
	"errors"
	"time"
)

var (
	ErrPriceIsRequired = errors.New("price is required")
	ErrInvalidPrice = errors.New("invalid price")
)

type Product struct {
	ID int64 `json:"id" gorm:"primaryKey;autoIncrement"`
	Name string `json:"name"`
	Price float64 `json:"price"`
	Active bool `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewProduct(name string, price float64, active bool) (*Product, error) {
	product := &Product{
		Name: name,
		Price: price,
		Active: active,
		CreatedAt: time.Now(),
	}
	err := product.ValidateFields()
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (p *Product) ValidateFields() error {
	if p.Name == "" {
		return ErrNameIsRequired
	}
	if p.Price == 0 {
		return ErrPriceIsRequired
	}
	if p.Price < 0 {
		return ErrInvalidPrice
	}
	return nil
}