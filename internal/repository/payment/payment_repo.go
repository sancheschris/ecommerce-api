package payment

import (
	"github.com/sancheschris/ecommerce-api/internal/model"
	"gorm.io/gorm"
)

type Payment struct {
	DB *gorm.DB
}

func NewPayment(db *gorm.DB) *Payment {
	return &Payment{DB: db}
}

func (p *Payment) Create(payment *model.Payment) error {
	return p.DB.Create(payment).Error
}

func (p *Payment) GetByID(id int) (*model.Payment, error) {
	var payment model.Payment
	err := p.DB.Preload("Order").First(&payment, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &payment, err
}

func (p *Payment) Update(payment *model.Payment) error {
	_, err := p.GetByID(payment.ID)
	if err != nil {
		return err
	}
	return p.DB.Save(payment).Error
}

func (p *Payment) Delete(id int) error {
	payment, err := p.GetByID(id)
	if err != nil {
		return err
	}
	return p.DB.Delete(payment).Error
}

func (p *Payment) GetByOrderID(orderID int) (*model.Payment, error) {
	var payment model.Payment
	err := p.DB.Preload("Order").First(&payment, "order_id = ?", orderID).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (p *Payment) GetByUserID(userID int) (*model.Payment, error) {
	var payment model.Payment
	err := p.DB.Preload("Order").Joins("JOIN orders ON payments.order_id = orders.id").
		Where("orders.user_id = ?", userID).
		First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}