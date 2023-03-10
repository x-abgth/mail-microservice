package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	// specify who is allowed to connect
	// middleware handler that enables CORS
	mux.Use(cors.Handler(cors.Options{
		// allowed any request start with https and http
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "PUT", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// to make sure that the service is responding
	mux.Use(middleware.Heartbeat("/ping"))

	// create a new route
	mux.Post("/authenticate", app.Authenticate)

	mux.Post("/register", app.Register)
	mux.Post("/get-user", app.GetUser)

	return mux
}
