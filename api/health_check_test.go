package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/freshusername/news-api/models"
)

func TestHealthCheckHandler(t *testing.T) {
	// Create an instance of the Application with the mock DB
	mockDB := &MockDatabaseRepo{
		HealthcheckFunc: func() (*models.Post, error) {
			// Return a healthy response
			return &models.Post{}, nil
		},
	}
	app := &Application{DB: mockDB}

	// Set up the request and recorder
	req, err := http.NewRequest(http.MethodGet, "/health-check", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	// Call the handler
	handler := http.HandlerFunc(app.HealthCheck)
	handler.ServeHTTP(rr, req)

	// Check the status code and response body
	if rr.Code != http.StatusOK {
		t.Errorf("expected status OK; got %v", rr.Code)
	}

	// Decode the response body to check the contents
	var gotPayload struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Version string `json:"version"`
	}
	err = json.NewDecoder(rr.Body).Decode(&gotPayload)
	if err != nil {
		t.Fatal(err)
	}
	wantPayload := struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Version string `json:"version"`
	}{
		Status:  "Active",
		Message: "News-api is Healthy",
		Version: "1.0.0",
	}
	if gotPayload != wantPayload {
		t.Errorf("expected payload %+v; got %+v", wantPayload, gotPayload)
	}
}
