package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/freshusername/news-api/models"
)

// Assuming MockDatabaseRepo is already defined and includes a dynamic CreatePost method.

func TestHandleCreatePost_Success(t *testing.T) {
	// Arrange
	mockDB := &MockDatabaseRepo{
		CreatePostFunc: func(post *models.Post) (*models.Post, error) {
			// Simulate successful creation by returning the same post with an ID
			post.ID = 1 // Assuming an ID is assigned by the database
			return post, nil
		},
	}
	app := &Application{DB: mockDB}

	postData := []byte(`{"title":"Test Title", "content":"Test Content"}`)
	req, _ := http.NewRequest("POST", "/posts", bytes.NewReader(postData))
	rr := httptest.NewRecorder()

	// Act
	app.HandleCreatePost(rr, req)

	// Assert
	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status %d; got %d", http.StatusCreated, rr.Code)
	}

	var createdPost models.Post
	if err := json.NewDecoder(rr.Body).Decode(&createdPost); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if createdPost.ID != 1 || createdPost.Title != "Test Title" || createdPost.Content != "Test Content" {
		t.Errorf("Unexpected post data returned: %+v", createdPost)
	}
}

func TestHandleCreatePost_Failure_ValidationError(t *testing.T) {
	// Arrange
	app := &Application{DB: &MockDatabaseRepo{}}

	// Invalid post data (empty title and content)
	postData := []byte(`{"title":"", "content":""}`)
	req, _ := http.NewRequest("POST", "/posts", bytes.NewReader(postData))
	rr := httptest.NewRecorder()

	// Act
	app.HandleCreatePost(rr, req)

	// Assert
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d; got %d", http.StatusBadRequest, rr.Code)
	}
}

// You can add more tests to simulate database errors or other validation errors.
func TestHandleCreatePost_Failure_DatabaseError(t *testing.T) {
	// Arrange
	mockDB := &MockDatabaseRepo{
		CreatePostFunc: func(post *models.Post) (*models.Post, error) {
			// Simulate a database error
			return nil, errors.New("database error")
		},
	}
	app := &Application{DB: mockDB}

	postData := []byte(`{"title":"Valid Title", "content":"Valid Content"}`)
	req, _ := http.NewRequest("POST", "/posts", bytes.NewReader(postData))
	rr := httptest.NewRecorder()

	// Act
	app.HandleCreatePost(rr, req)

	// Assert
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d; got %d", http.StatusInternalServerError, rr.Code)
	}
}

func TestHandleCreatePost_Failure_MissingField(t *testing.T) {
	// Arrange
	app := &Application{DB: &MockDatabaseRepo{}}

	// Missing "content" field
	postData := []byte(`{"title":"Title Without Content"}`)
	req, _ := http.NewRequest("POST", "/posts", bytes.NewReader(postData))
	rr := httptest.NewRecorder()

	// Act
	app.HandleCreatePost(rr, req)

	// Assert
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d; got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestHandleCreatePost_Success_WithValidation(t *testing.T) {
	// Arrange
	mockDB := &MockDatabaseRepo{
		CreatePostFunc: func(post *models.Post) (*models.Post, error) {
			// Assuming an ID is assigned by the database and all validations pass
			post.ID = 2
			return post, nil
		},
	}
	app := &Application{DB: mockDB}

	postData := []byte(`{"title":"Another Test Title", "content":"Another Test Content"}`)
	req, _ := http.NewRequest("POST", "/posts", bytes.NewReader(postData))
	rr := httptest.NewRecorder()

	// Act
	app.HandleCreatePost(rr, req)

	// Assert
	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status %d; got %d", http.StatusCreated, rr.Code)
	}

	var createdPost models.Post
	if err := json.NewDecoder(rr.Body).Decode(&createdPost); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if createdPost.ID != 2 || createdPost.Title != "Another Test Title" || createdPost.Content != "Another Test Content" {
		t.Errorf("Unexpected post data returned: %+v", createdPost)
	}
}
