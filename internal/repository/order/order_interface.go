package order

import "github.com/sancheschris/ecommerce-api/internal/model"

type OrderInterface interface {
	GetOrders() ([]model.Order, error)
	GetOrderByID(id int) (*model.Order, error)
	GetOrdersByUserID(userID int64) ([]model.Order, error)
	CreateOrder(order *model.Order) error
	UpdateOrder(order *model.Order) error
	DeleteOrder(id int64) error

	AddOrderItem(orderID int64, item *model.OrderItem) error
	UpdateOrderItem(orderID int64, item *model.OrderItem) error
	RemoveOrderItem(orderID int64, itemID int64) error
	GetOrderItems(orderID int64) ([]model.Order, error)
}