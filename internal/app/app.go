package app

import (
	"fmt"
	"net/http"
	"shop/internal/env"
	"shop/internal/repositories"
	"shop/internal/services"
	"time"
)

type Application struct {
	port      int
	jwtSecret string

	userRepo     repositories.UserRepository
	productRepo  repositories.ProductRepository
	cartRepo     repositories.CartRepository
	cartItemRepo repositories.CartItemRepository

	userService    *services.UserService
	cartService    *services.CartService
	productService *services.ProductService
}

func GetApplication(
	userRepo repositories.UserRepository,
	productRepo repositories.ProductRepository,
	cartRepo repositories.CartRepository,
	cartItemRepo repositories.CartItemRepository) *Application {

	return &Application{
		port:      env.GetEnvInt("PORT", 8080),
		jwtSecret: env.GetEnvString("JWT_SECRET", "some_secret"),

		userRepo:     userRepo,
		productRepo:  productRepo,
		cartRepo:     cartRepo,
		cartItemRepo: cartItemRepo,

		userService:    services.NewUserService(userRepo, cartRepo, env.GetEnvString("JWT_SECRET", "some_secret")),
		cartService:    services.NewCartService(cartRepo, cartItemRepo),
		productService: services.NewProductService(productRepo),
	}
}

func (app *Application) Serve() error {
	r := GetRouter(app)

	server := http.Server{
		Addr:         fmt.Sprintf(":%d", app.port),
		Handler:      r.Route(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server.ListenAndServe()
}
