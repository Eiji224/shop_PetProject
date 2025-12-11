package main

import (
	"fmt"
	"log"
	_ "shop/docs"
	"shop/internal/app"
	"shop/internal/env"
	"shop/internal/repositories"
	"time"

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

	var db *gorm.DB
	var err error
	maxRetries := 10
	retryDelay := 3 * time.Second

	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Println("Успешно подключено к базе данных")
			return db
		}

		log.Printf("Попытка подключения к БД %d/%d не удалась: %v", i+1, maxRetries, err)
		if i < maxRetries-1 {
			log.Printf("Повторная попытка через %v...", retryDelay)
			time.Sleep(retryDelay)
		}
	}

	log.Printf("Не удалось подключиться к БД после %d попыток. DSN: %s", maxRetries, dsn)
	panic(fmt.Sprintf("не удалось подключиться к базе данных: %v", err))
}
