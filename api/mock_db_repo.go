package main

import (
	"database/sql"

	"github.com/freshusername/news-api/models"
)

type MockDatabaseRepo struct {
	HealthcheckFunc func() (*models.Post, error)
	CreatePostFunc  func(post *models.Post) (*models.Post, error)
}

func (m *MockDatabaseRepo) Connection() *sql.DB {
	// Mock the Connection method if necessary
	return nil
}

func (m *MockDatabaseRepo) Healthcheck() (*models.Post, error) {
	return m.HealthcheckFunc()
}

// You must add mock implementations for all other methods defined in the DatabaseRepo interface
func (m *MockDatabaseRepo) GetAllPosts() ([]*models.Post, error) {
	// Return a slice of posts or an error based on your test needs
	return nil, nil
}

func (m *MockDatabaseRepo) CreatePost(post *models.Post) (*models.Post, error) {
	// Return a new post or an error based on your test needs
	return nil, nil
}

func (m *MockDatabaseRepo) UpdatePost(id int32, item *models.Post) (*models.Post, error) {
	// Return an updated post or an error based on your test needs
	return nil, nil
}

func (m *MockDatabaseRepo) DeletePost(id int32) (int32, error) {
	// Return the ID of the deleted post or an error based on your test needs
	return 0, nil
}
