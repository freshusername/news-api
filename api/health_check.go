package main

import (
	"fmt"
	"net/http"
)

func (app *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	var payload = struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Version string `json:"version"`
	}{
		Status:  "Active",
		Message: "News-api is Healthy",
		Version: "1.0.0",
	}

	check, err := app.DB.Healthcheck()

	if check != nil {
		_ = app.writeJSON(w, http.StatusOK, payload)
	} else {
		_ = app.writeJSON(w, http.StatusOK, fmt.Sprintf("Unhealthy. Error: %s", err))
	}
}
