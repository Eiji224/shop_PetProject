package handlers

import (
	"net/http"
	"shop/internal/database"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Router struct {
	jwtSecret string
	models    *database.Models
}

func GetRouter(jwtSecret string, models *database.Models) *Router {
	return &Router{jwtSecret: jwtSecret, models: models}
}

func (r *Router) Route() http.Handler {
	g := gin.Default()

	v1 := g.Group("/api/v1")
	{
		v1.POST("/register", r.Register)
		v1.POST("/login", r.Login)
	}

	g.GET("/swagger/*any", func(c *gin.Context) {
		if c.Request.RequestURI == "/swagger/" {
			c.Redirect(http.StatusFound, "/swagger/index.html")
		}
		ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("http://localhost:8080/swagger/doc.json"))(c)
	})

	return g
}
