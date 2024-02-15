package main

import (
	"net/http"
)

func (app *application) HandleSwagger(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./openapi/swagger.json")
}
