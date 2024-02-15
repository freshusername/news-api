package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWriteJSON(t *testing.T) {
	// Create an instance of your Application struct with necessary fields initialized.
	app := &Application{}

	// Create a ResponseRecorder to record the response.
	rr := httptest.NewRecorder()

	// The data that we want to write as JSON.
	testData := struct {
		Name string `json:"name"`
	}{
		Name: "Test",
	}

	// Call writeJSON.
	err := app.writeJSON(rr, http.StatusOK, testData)
	if err != nil {
		t.Errorf("writeJSON returned an error: %v", err)
	}

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("writeJSON returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response content type is JSON.
	expectedContentType := "Application/json"
	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("writeJSON returned wrong Content-Type header: got %v want %v", contentType, expectedContentType)
	}

	// Check the response body is what we expect.
	expectedBody := `{"name":"Test"}`
	if rr.Body.String() != expectedBody {
		t.Errorf("writeJSON returned unexpected body: got %v want %v", rr.Body.String(), expectedBody)
	}
}

func TestReadJSON(t *testing.T) {
	// Create an instance of your Application struct with necessary fields initialized.
	app := &Application{}

	// Simulate a JSON payload to read.
	jsonPayload := `{"name":"Test"}`
	req := httptest.NewRequest("POST", "/url", strings.NewReader(jsonPayload))

	// Create a ResponseRecorder to record the response, needed for MaxBytesReader.
	rr := httptest.NewRecorder()

	// The struct that we expect to be populated.
	testData := struct {
		Name string `json:"name"`
	}{}

	// Call readJSON.
	err := app.readJSON(rr, req, &testData)
	if err != nil {
		t.Fatalf("readJSON returned an error: %v", err)
	}

	// Check if the payload was correctly decoded.
	expectedName := "Test"
	if testData.Name != expectedName {
		t.Errorf("readJSON did not decode correctly: got %v want %v", testData.Name, expectedName)
	}
}

func TestErrorJSON(t *testing.T) {
	// Create an instance of your Application struct with necessary fields initialized.
	app := &Application{}

	// Create a ResponseRecorder to record the response.
	rr := httptest.NewRecorder()

	// Call errorJSON with a sample error.
	testError := errors.New("this is a test error")
	err := app.errorJSON(rr, testError)
	if err != nil {
		t.Errorf("errorJSON returned an error: %v", err)
	}

	// Check the status code is what we expect.
	expectedStatusCode := http.StatusBadRequest
	if status := rr.Code; status != expectedStatusCode {
		t.Errorf("errorJSON returned wrong status code: got %v want %v", status, expectedStatusCode)
	}

	// Check the response body is what we expect.
	expectedBody := `{"error":true,"message":"this is a test error"}`
	if rr.Body.String() != expectedBody {
		t.Errorf("errorJSON returned unexpected body: got %v want %v", rr.Body.String(), expectedBody)
	}
}
