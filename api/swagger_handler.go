package main

import (
	"net/http"
)

func (app *Application) HandleSwagger(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./openapi/swagger.json")
}
