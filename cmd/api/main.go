package main

import (
	"fmt"
	_ "shop/docs"
	"shop/internal/app"
	"shop/internal/env"
	"shop/internal/repositories"

	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// @title Shop
// @version 1.0
// @description A Shop wrote by Go using Gin framework
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	db := connectDB()
	userRep := repositories.NewUserRepository(db)
	productRep := repositories.NewProductRepository(db)
	cartRep := repositories.NewCartRepository(db)
	cartItemRep := repositories.NewCartItemRepository(db)

	application := app.GetApplication(userRep, productRep, cartRep, cartItemRep)

	if err := application.Serve(); err != nil {
		panic(err)
	}
}

func connectDB() *gorm.DB {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		env.GetEnvString("DB_USER", "mysql"),
		env.GetEnvString("DB_PASSWORD", "mysql"),
		env.GetEnvString("DB_HOST", "localhost"),
		env.GetEnvString("DB_PORT", "3306"),
		env.GetEnvString("DB_NAME", "mysql"),
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}
