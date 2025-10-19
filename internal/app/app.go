package app

import (
	"fmt"
	"net/http"
	"shop/internal/database"
	"shop/internal/env"
	"shop/internal/handlers"
	"shop/internal/services"
	"time"
)

type Application struct {
	port      int
	jwtSecret string
	models    database.Models
	services  services.Services
}

func GetApplication(models database.Models) *Application {
	services := services.NewServices(env.GetEnvString("JWT_SECRET", "some_secret"), &models)

	return &Application{
		port:      env.GetEnvInt("PORT", 8080),
		jwtSecret: env.GetEnvString("JWT_SECRET", "some_secret"),
		models:    models,
		services:  services,
	}
}

func (app *Application) Serve() error {
	r := handlers.GetRouter(app.jwtSecret, &app.models, &app.services)

	server := http.Server{
		Addr:         fmt.Sprintf(":%d", app.port),
		Handler:      r.Route(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server.ListenAndServe()
}
