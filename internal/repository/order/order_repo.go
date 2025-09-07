package order

import (
	"fmt"

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
	err := o.DB.Preload("User").Preload("Items").Preload("Items.Product").Preload("Payments").Preload("Payments.Order").Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, err
}

func (o *Order) GetOrderByID(id int) (*model.Order, error) {
	var order model.Order
	err := o.DB.Preload("Items.Product").Preload("Payments").First(&order, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &order, err
}

func (o *Order) UpdateOrder(order *model.Order) error {
	if err := o.DB.Model(order).Updates(order).Error; err != nil {
		return err
	}
	if err := o.DB.Model(order).Association("Items").Replace(order.Items); err != nil {
		return err
	}
	if err := o.DB.Model(order).Association("Payments").Replace(order.Payments); err != nil {
		return err
	}
	return nil
}

func (o *Order) DeleteOrder(id int) error {
	var order []model.Order
	err := o.DB.First(&order).Error
	if err != nil {
		return err
	}
	return o.DB.Delete(&order, "id = ?", id).Error
}

func (o *Order) GetOrdersByUserID(userID int) ([]model.Order, error) {
	var orders []model.Order
	err := o.DB.Where("user_id = ?", userID).Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (o *Order) AddOrderItem(orderID int, item *model.OrderItem) error {
	item.OrderID = orderID
	return o.DB.Create(&item).Error
}

func (o *Order) UpdateOrderItem(orderID int, item *model.OrderItem) error {
	_, err := o.GetOrderByID(orderID)
	if err != nil {
		return err
	}
	if item.OrderID != orderID {
		return fmt.Errorf("Order items does not belong to order")
	}
	return o.DB.Save(item).Error
}

func (o *Order) GetOrderItems(orderID int) ([]model.OrderItem, error) {
	var orderItems []model.OrderItem
	err := o.DB.Find(&orderItems, "order_id = ?", orderID).Error
	if err != nil {
		return nil, err
	}
	return orderItems, nil
}

func (o *Order) RemoveOrderItem(orderID int, itemID int) error {
    result := o.DB.Delete(&model.OrderItem{}, "order_id = ? AND id = ?", orderID, itemID)
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return fmt.Errorf("no order item found to delete")
    }
    return nil
}