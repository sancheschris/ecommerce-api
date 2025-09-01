package order

import (
	"github.com/sancheschris/ecommerce-api/internal/model"
	"gorm.io/gorm"
)

type Order struct {
	DB *gorm.DB
}

func NewOrder(db *gorm.DB) *Order {
	return &Order{DB: db}
}

func (o *Order) CreateOrder(order *model.Order) error {
	return o.DB.Create(order).Error
}

func (o *Order) GetOrders() ([]model.Order, error) {
	var orders []model.Order
	err := o.DB.Find(orders).Error
	if err != nil {
		return nil, err
	}
	return orders, err
}

func (o *Order) GetOrderByID(id int64) (*model.Order, error) {
	var order model.Order
	err := o.DB.First(&order, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &order, err
}

func (o *Order) UpdateOrder(order *model.Order) error {
	_, err := o.GetOrderByID(order.ID)
	if err != nil {
		return err
	}
	return o.DB.Save(order).Error
}