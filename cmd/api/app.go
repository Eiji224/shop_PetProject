package main

import (
	"fmt"
	"net/http"
	"shop/internal/database"
	"time"
)

type application struct {
	port      int
	jwtSecret string
	models    database.Models
}

func (app *application) serve() error {
	server := http.Server{
		Addr:         fmt.Sprintf(":%d", app.port),
		Handler:      app.route(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server.ListenAndServe()
}
