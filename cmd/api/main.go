package main

import (
	"fmt"
	"shop/internal/app"
	"shop/internal/database"
	"shop/internal/env"

	_ "shop/docs"

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
	err := db.AutoMigrate(&database.User{}, &database.Category{}, &database.Product{}, &database.Cart{},
		&database.CartItem{}, &database.Order{}, &database.OrderItem{})
	if err != nil {
		panic(err)
	}

	models := database.NewModels(db)
	application := app.GetApplication(models)

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
		env.GetEnvString("DB_PORT", "5432"),
		env.GetEnvString("DB_NAME", "mysql"),
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}
