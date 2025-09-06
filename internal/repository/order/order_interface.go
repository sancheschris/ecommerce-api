package order

import "github.com/sancheschris/ecommerce-api/internal/model"

type OrderInterface interface {
	GetOrders() ([]model.Order, error)
	GetOrderByID(id int) (*model.Order, error)
	GetOrdersByUserID(userID int) ([]model.Order, error)
	CreateOrder(order *model.Order) error
	UpdateOrder(order *model.Order) error
	DeleteOrder(id int) error

	AddOrderItem(orderID int, item *model.OrderItem) error
	UpdateOrderItem(orderID int, item *model.OrderItem) error
	RemoveOrderItem(orderID int, itemID int) error
	GetOrderItems(orderID int) ([]model.OrderItem, error)
}