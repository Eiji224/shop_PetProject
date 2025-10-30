package app

import (
	"net/http"
	"shop/internal/handlers"
	"shop/internal/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Router struct {
	jwtSecret string

	userHandler    *handlers.UserHandler
	productHandler *handlers.ProductHandler
	cartHandler    *handlers.CartHandler

	middleware *middleware.Middleware
}

func GetRouter(app *Application) *Router {
	return &Router{
		jwtSecret: app.jwtSecret,

		userHandler:    handlers.NewUserHandler(app.userService),
		productHandler: handlers.NewProductHandler(app.productRepo, app.userService, app.productService),
		cartHandler:    handlers.NewCartHandler(app.cartRepo, app.cartItemRepo, app.productRepo, app.userService, app.cartService),

		middleware: middleware.GetMiddleware(app.jwtSecret, app.userRepo),
	}
}

func (r *Router) Route() http.Handler {
	g := gin.Default()

	v1 := g.Group("/api/v1")
	{
		v1.GET("/products", r.productHandler.GetAllProducts)
		v1.GET("/products/:id", r.productHandler.GetProduct)
		v1.GET("/users/:id/products", r.productHandler.GetProductsBySeller)

		v1.POST("/auth/register", r.userHandler.Register)
		v1.POST("/auth/login", r.userHandler.Login)
	}

	authGroup := v1.Group("/")
	authGroup.Use(r.middleware.AuthMiddleware())
	{
		authGroup.POST("/products", r.productHandler.CreateProduct)
		authGroup.PUT("/products/:id", r.productHandler.UpdateProduct)
		authGroup.DELETE("/products/:id", r.productHandler.DeleteProduct)

		authGroup.GET("/cart/item", r.cartHandler.GetCartItems)
		authGroup.POST("/cart/item", r.cartHandler.AddCartItem)
		authGroup.PATCH("/cart/item/:id", r.cartHandler.UpdateCartItemQuantity)
		authGroup.DELETE("/cart/item/:id", r.cartHandler.DeleteCartItem)
		authGroup.DELETE("/cart/item", r.cartHandler.DeleteAllCartItems)
	}

	g.GET("/swagger/*any", func(c *gin.Context) {
		if c.Request.RequestURI == "/swagger/" {
			c.Redirect(http.StatusFound, "/swagger/index.html")
		}
		ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("http://localhost:8080/swagger/doc.json"))(c)
	})

	return g
}
