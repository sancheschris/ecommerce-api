package main

import (
	"github.com/sancheschris/ecommerce-api/configs"
	"github.com/sancheschris/ecommerce-api/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	_, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(sqlite.Open("ecommerce-api.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&model.User{}, &model.Product{}, &model.Payment{}, &model.Order{}, &model.OrderItem{})
}