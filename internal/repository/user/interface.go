package repository

import "github.com/sancheschris/ecommerce-api/internal/model"

type UserInterface interface {
	Create(user *model.User) error
	FindByEmail(emailId string) (*model.User, error)
}