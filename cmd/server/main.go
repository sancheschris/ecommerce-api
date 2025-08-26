package main

import (
	"github.com/sancheschris/ecommerce-api/configs"
	"github.com/sancheschris/ecommerce-api/internal/entity"
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
	db.AutoMigrate(&entity.User{}, &entity.Product{}, &entity.Payment{}, &entity.Order{}, &entity.OrderItem{})
}