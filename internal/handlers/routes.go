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
		v1.POST("/auth/register", r.Register)
		v1.POST("/auth/login", r.Login)
	}

	authGroup := v1.Group("/")
	authGroup.Use(m.AuthMiddleware())
	{
		authGroup.GET("/products", r.GetAllProducts)
	}

	g.GET("/swagger/*any", func(c *gin.Context) {
		if c.Request.RequestURI == "/swagger/" {
			c.Redirect(http.StatusFound, "/swagger/index.html")
		}
		ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("http://localhost:8080/swagger/doc.json"))(c)
	})

	return g
}
