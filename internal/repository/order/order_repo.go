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
	err := o.DB.
		Preload("User").
		Preload("Items.Product").
		Preload("Payments").
		Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, err
}

func (o *Order) GetOrderByID(id int) (*model.Order, error) {
	var order model.Order
	err := o.DB.
		Preload("User").
		Preload("Items.Product").
		Preload("Payments").
		First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, err
}

func (o *Order) UpdateOrder(order *model.Order) error {
    return o.DB.Transaction(func(tx *gorm.DB) error {
        // update only allowed fields (to avoid zero-value issues)
        if err := tx.Model(&model.Order{ID: order.ID}).Updates(map[string]any{
            "status":      order.Status,
            "total_price": order.TotalPrice,
            "currency":    order.Currency,
        }).Error; err != nil {
            return err
        }

        // replace items/payments if that's desired
        for i := range order.Items {
            order.Items[i].OrderID = order.ID
        }
        if err := tx.Model(&order).Association("Items").Replace(order.Items); err != nil {
            return err
        }

        for i := range order.Payments {
            order.Payments[i].OrderID = order.ID
        }
        if err := tx.Model(&order).Association("Payments").Replace(order.Payments); err != nil {
            return err
        }

        return nil
    })
}

func (o *Order) DeleteOrder(id int) error {
    var order model.Order
    if err := o.DB.First(&order, id).Error; err != nil {
        return err
    }
    return o.DB.Delete(&order).Error
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
	return o.DB.Create(item).Error
}

func (o *Order) UpdateOrderItem(orderID int, item *model.OrderItem) error {
    if item.OrderID != orderID {
        return fmt.Errorf("order item does not belong to order %d", orderID)
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