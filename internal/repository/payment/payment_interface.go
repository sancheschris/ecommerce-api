package payment

import "github.com/sancheschris/ecommerce-api/internal/model"

type PaymentInterface interface{
	Create(payment *model.Payment) error
	GetByID(id int) (*model.Payment, error)
	Update(payment *model.Payment) error
	Delete(id int) error

	GetByOrderID(orderID int) (*model.Payment, error)
	GetByUserID(userID int) (*model.Payment, error)
	GetByStatus(status string) ([]*model.Payment, error)
}