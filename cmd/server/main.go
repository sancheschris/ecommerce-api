package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sancheschris/ecommerce-api/configs"
	"github.com/sancheschris/ecommerce-api/internal/handler"
	"github.com/sancheschris/ecommerce-api/internal/model"
	orderRepo "github.com/sancheschris/ecommerce-api/internal/repository/order"
	productRepo "github.com/sancheschris/ecommerce-api/internal/repository/product"
	userRepo "github.com/sancheschris/ecommerce-api/internal/repository/user"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(sqlite.Open("ecommerce.db"), &gorm.Config{})
	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
    // config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName)
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&model.User{}, &model.Product{}, &model.Payment{}, &model.Order{}, &model.OrderItem{})

	userDB := userRepo.NewUser(db)
	userHandler := handler.NewUserHandler(userDB)

	productDB := productRepo.NewProduct(db)
	productHandler := handler.NewProductHandler(productDB)

	orderDB := orderRepo.NewOrder(db)
	orderHandler := handler.NewOrderHandler(orderDB)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.WithValue("jwt", configs.TokenAuth))
	r.Use(middleware.WithValue("JwtExpiresIn", configs.JwtExpiresIn))
	
	r.Post("/users", userHandler.Create)
	r.Post("/users/generate_token", userHandler.GetJWT)
	r.Get("/users/orders", userHandler.GetOrders)

	r.Post("/products", productHandler.Create)
	r.Get("/products/{id}", productHandler.GetProductByID)
	r.Get("/products", productHandler.GetProducts)
	r.Put("/products/{id}", productHandler.UpdateProduct)
	r.Delete("/products/{id}", productHandler.DeleteProduct)

	r.Post("/orders", orderHandler.CreateOrder)
	r.Get("/orders", orderHandler.GetOrders)
	r.Get("/orders/{id}", orderHandler.GetOrderByID)
	r.Put("/orders/{id}", orderHandler.UpdateOrder)

	http.ListenAndServe(":8080", r)
}