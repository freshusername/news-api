package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *Application) routes() http.Handler {
	// create a router mux
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Get("/", app.HealthCheck)
	mux.Get("/posts", app.HandleGetPosts)
	mux.Post("/posts", app.HandleCreatePost)
	mux.Put("/posts/{id}", app.HandleUpdatePost)
	mux.Delete("/posts/{id}", app.HandleDeletePost)

	//openapi specification
	mux.Get("/swagger", app.HandleSwagger)

	return mux
}
