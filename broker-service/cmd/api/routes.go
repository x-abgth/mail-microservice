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

	/*
		The middleware.Heartbeat function installs a handler on the router that responds to requests
		at the specified path with a simple "pong" message. This allows clients to check the status of
		the server by making a request to the endpoint.
	*/
	mux.Use(middleware.Heartbeat("/ping"))

	// handle function
	mux.Post("/", app.Broker)
	// this entrypoint of the where
	// every microservice will request in
	mux.Post("/handle", app.HanldeSubmission)

	return mux
}
