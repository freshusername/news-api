package database

import (
	"database/sql"

	"github.com/freshusername/news-api/models"
)

type DatabaseRepo interface {
	Connection() *sql.DB
	Healthcheck() (*models.Post, error)
}
