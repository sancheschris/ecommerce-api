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
	err := o.DB.Preload("Items").Preload("Payments").Find(&orders).Error
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

func (o *Order) DeleteOrder(id int64) error {
	var order []model.Order
	err := o.DB.First(&order).Error
	if err != nil {
		return err
	}
	return o.DB.Delete(&order, "id = ?", id).Error
}

func (o *Order) GetOrdersByUserID(userID int64) ([]model.Order, error) {
	var orders []model.Order
	err := o.DB.Where("user_id = ?", userID).Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (o *Order) AddOrderItem(orderID int64, item *model.OrderItem) error {
	item.OrderID = orderID
	return o.DB.Create(item).Error
}