package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/freshusername/news-api/models"
)

type PostgresDBRepo struct {
	DB *sql.DB
}

const dbTimeout = time.Second * 3

func (m *PostgresDBRepo) Connection() *sql.DB {
	return m.DB
}

func (m *PostgresDBRepo) GetAllPosts() ([]*models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT id, title, content, created_at, updated_at
		FROM public.posts
		ORDER BY created_at DESC
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post

	for rows.Next() {
		var post models.Post

		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, err
		}

		posts = append(posts, &post)
	}

	return posts, nil
}

func (repo *PostgresDBRepo) CreatePost(post *models.Post) (*models.Post, error) {
	query := `
        INSERT INTO public.posts (title, content, created_at, updated_at) 
        VALUES ($1, $2, NOW(), NOW())
        RETURNING id, title, content, created_at, updated_at
    `

	row := repo.DB.QueryRow(query, post.Title, post.Content)

	newPost := &models.Post{}
	err := row.Scan(&newPost.ID, &newPost.Title, &newPost.Content, &newPost.CreatedAt, &newPost.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return newPost, nil
}

func (m *PostgresDBRepo) Healthcheck() (*models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT id, title, content, created_at, updated_at
		FROM public.posts
		ORDER BY created_at DESC
		LIMIT 1
	`

	row := m.DB.QueryRowContext(ctx, query)

	post := &models.Post{}
	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (m *PostgresDBRepo) UpdatePost(id int32, post *models.Post) (*models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		UPDATE public.posts
		SET title = $2, content = $3, updated_at = NOW()
		WHERE id = $1
		RETURNING id, title, content, created_at, updated_at
    `

	row := m.DB.QueryRowContext(ctx, query, id, post.Title, post.Content)

	updatedPost := &models.Post{}

	err := row.Scan(&updatedPost.ID, &updatedPost.Title, &updatedPost.Content, &updatedPost.CreatedAt, &updatedPost.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("no rows were updated, post may not exist")
		}
		return nil, err
	}

	return updatedPost, nil
}

func (m *PostgresDBRepo) DeletePost(id int32) (int32, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `DELETE FROM public.posts WHERE id = $1`

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	if rowsAffected == 0 {
		return 0, errors.New("no rows were deleted, post may not exist")
	}

	return id, nil
}
