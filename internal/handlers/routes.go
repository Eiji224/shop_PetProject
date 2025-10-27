package handlers

import (
	"net/http"
	"shop/internal/database"
	"shop/internal/middleware"
	"shop/internal/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Router struct {
	jwtSecret string
	models    *database.Models
	services  *services.Services
}

func GetRouter(jwtSecret string, models *database.Models, services *services.Services) *Router {
	return &Router{jwtSecret: jwtSecret, models: models, services: services}
}

func (r *Router) Route() http.Handler {
	g := gin.Default()
	m := middleware.GetMiddleware(r.jwtSecret, r.models)

	v1 := g.Group("/api/v1")
	{
		v1.GET("/products", r.getAllProducts)
		v1.GET("/products/:id", r.getProduct)
		v1.GET("/users/:id/products", r.getProductsBySeller)

		v1.POST("/auth/register", r.Register)
		v1.POST("/auth/login", r.Login)
	}

	authGroup := v1.Group("/")
	authGroup.Use(m.AuthMiddleware())
	{
		authGroup.POST("/products", r.createProduct)
		authGroup.PUT("/products/:id", r.updateProduct)
		authGroup.DELETE("/products/:id", r.deleteProduct)

		authGroup.GET("/cart/item", r.getCartItems)
		authGroup.POST("/cart/item", r.addCartItem)
		authGroup.PATCH("/cart/item/:id", r.updateCartItemQuantity)
		authGroup.DELETE("/cart/item/:id", r.deleteCartItem)
		authGroup.DELETE("/cart/item", r.deleteAllCartItems)
	}

	g.GET("/swagger/*any", func(c *gin.Context) {
		if c.Request.RequestURI == "/swagger/" {
			c.Redirect(http.StatusFound, "/swagger/index.html")
		}
		ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("http://localhost:8080/swagger/doc.json"))(c)
	})

	return g
}
