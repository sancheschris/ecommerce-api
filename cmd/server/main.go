package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/sancheschris/ecommerce-api/configs"
	_ "github.com/sancheschris/ecommerce-api/docs"
	"github.com/sancheschris/ecommerce-api/internal/handler"
	"github.com/sancheschris/ecommerce-api/internal/model"
	orderRepo "github.com/sancheschris/ecommerce-api/internal/repository/order"
	paymentRepo "github.com/sancheschris/ecommerce-api/internal/repository/payment"
	productRepo "github.com/sancheschris/ecommerce-api/internal/repository/product"
	userRepo "github.com/sancheschris/ecommerce-api/internal/repository/user"
	"github.com/sancheschris/ecommerce-api/internal/service"
	"github.com/sancheschris/ecommerce-api/pkg/payment"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// @title           Ecommerce API
// @version         1.0
// @description     Ecommerce API
// @termsOfService  http://swagger.io/terms/

// @contact.name   Christian Santos
// @contact.url    https://github.com/sancheschris
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /
// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization
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

	if configs.StripeSecretKey == "" {
        log.Fatal("STRIPE_SECRET_KEY configuration is required")
    }

	stripeClient := payment.NewStripeClient(configs.StripeSecretKey)

	paymentRepo := paymentRepo.NewPayment(db)
	paymentService := service.NewPaymentService(paymentRepo, stripeClient)
	paymentHandler := handler.NewPaymentHandler(paymentService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.WithValue("jwt", configs.TokenAuth))
	r.Use(middleware.WithValue("JwtExpiresIn", configs.JwtExpiresIn))
	
	r.Post("/users", userHandler.Create)
	r.Post("/users/generate_token", userHandler.GetJWT)
	r.With(jwtauth.Verifier(configs.TokenAuth), jwtauth.Authenticator).Get("/users/orders", userHandler.GetOrders)

	r.Route("/products", func(r chi.Router) {
		r.Use(jwtauth.Verifier(configs.TokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Post("/", productHandler.Create)
		r.Get("/{id}", productHandler.GetProductByID)
		r.Get("/", productHandler.GetProducts)
		r.Put("/{id}", productHandler.UpdateProduct)
		r.Delete("/{id}", productHandler.DeleteProduct)
	})

	r.Route("/orders", func(r chi.Router) {
		r.Use(jwtauth.Verifier(configs.TokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Post("/", orderHandler.CreateOrder)
		r.Get("/", orderHandler.GetOrders)
		r.Get("/{id}", orderHandler.GetOrderByID)
		r.Get("/{id}/orders", orderHandler.GetOrdersByUserID)
		r.Put("/{id}", orderHandler.UpdateOrder)
		r.Delete("/{id}", orderHandler.DeleteOrder)
	})
	
	r.Route("/payments", func(r chi.Router) {
		r.Use(jwtauth.Verifier(configs.TokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Post("/", paymentHandler.CreatePayment)
		r.Get("/{id}/status", paymentHandler.GetPaymentStatus)
	})

	r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8080/docs/doc.json")))

	http.ListenAndServe(":8080", r)
}