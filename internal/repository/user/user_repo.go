package repository

import (
	"github.com/sancheschris/ecommerce-api/internal/model"
	"gorm.io/gorm"
)

type User struct {
	DB *gorm.DB
}

func NewUser(db *gorm.DB) *User {
	return &User{DB: db}
}

func (u *User) Create(user *model.User) error {
	return u.DB.Create(user).Error
}

func (u *User) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := u.DB.First(&user, "email = ?", email).Error
	return &user, err
}

func (u *User) GetOrders() ([]model.User, error) {
	var users []model.User
	err := u.DB.Preload("Orders").Find(&users).Error
	return users, err
}