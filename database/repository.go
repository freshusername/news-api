package database

import (
	"database/sql"

	"github.com/freshusername/news-api/models"
)

type DatabaseRepo interface {
	Connection() *sql.DB
	Healthcheck() (*models.Post, error)
	GetAllPosts() ([]*models.Post, error)
	CreatePost(item *models.Post) (*models.Post, error)
	UpdatePost(id int32, item *models.Post) (*models.Post, error)
	DeletePost(id int32) (int32, error)
}
