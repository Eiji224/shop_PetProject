package services

import "shop/internal/database"

type Services struct {
	AuthService *AuthService
}

func NewServices(jwtSecret string, models *database.Models) Services {
	return Services{
		&AuthService{jwtSecret: jwtSecret, models: models},
	}
}
