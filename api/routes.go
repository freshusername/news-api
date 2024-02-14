package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	// create a router mux
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	//mux.Use(app.enableCORS)

	mux.Get("/", app.HealthCheck)
	mux.Get("/news", app.HandleGetPosts)
	mux.Post("/news", app.HandleCreatePost)
	mux.Put("/news/{id}", app.HandleUpdatePost)
	mux.Delete("/news/{id}", app.HandleDeletePost)

	return mux
}
