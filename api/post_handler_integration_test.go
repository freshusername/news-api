package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"testing"

	_ "github.com/lib/pq"

	"github.com/freshusername/news-api/database"
	"github.com/freshusername/news-api/models"
	"github.com/go-chi/chi/v5"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	db      *sql.DB
	cleanup func()
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	var err error
	db, cleanup, err = setupPostgresContainer(ctx)
	if err != nil {
		log.Fatalf("Could not start postgres container: %s", err)
	}

	// Run the tests
	code := m.Run()

	// Cleanup resources after tests are done
	cleanup()

	os.Exit(code)
}

func setupPostgresContainer(ctx context.Context) (*sql.DB, func(), error) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:13",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       "testdb",
			"POSTGRES_USER":     "user",
			"POSTGRES_PASSWORD": "password",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}
	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, nil, err
	}

	// Closure to clean up the container
	cleanup := func() {
		postgresC.Terminate(ctx)
	}

	// Construct DSN
	port, _ := postgresC.MappedPort(ctx, "5432")
	dsn := fmt.Sprintf("host=localhost port=%s user=user password=password dbname=testdb sslmode=disable", port.Port())

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, cleanup, err
	}

	// Applying Goose migrations
	if err := applyGooseMigrations(dsn); err != nil {
		cleanup()
		log.Fatalf("Failed to apply goose migrations: %v", err)
	}
	return db, cleanup, nil
}

func TestHandleCreatePost_Integration(t *testing.T) {
	ctx := context.Background()
	db, cleanup, err := setupPostgresContainer(ctx)
	if err != nil {
		t.Fatalf("Could not start postgres container: %s", err)
	}
	defer cleanup()

	// Setup application with the connected database
	app := &Application{
		DB: &database.PostgresDBRepo{DB: db},
	}

	// Define the test case
	t.Run("Create Post Success", func(t *testing.T) {
		body := `{"title":"Integration Test Post", "content":"Content of the post"}`
		req := httptest.NewRequest("POST", "/posts", bytes.NewBufferString(body))
		rr := httptest.NewRecorder()

		app.HandleCreatePost(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusCreated)
		}

		var post models.Post
		if err := json.NewDecoder(rr.Body).Decode(&post); err != nil {
			t.Fatalf("Could not decode response: %v", err)
		}

		if post.Title != "Integration Test Post" {
			t.Errorf("Handler returned unexpected body: got title %v want %v", post.Title, "Integration Test Post")
		}
	})
}

func applyGooseMigrations(dsn string) error {
	migrationsDir := "../database/migrations"
	cmd := exec.Command("goose", "postgres", dsn, "up")
	cmd.Dir = migrationsDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("goose up failed: %v; output: %s", err, string(output))
	}
	return nil
}

func TestHandleGetPosts_Integration(t *testing.T) {
	ctx := context.Background()
	db, cleanup, err := setupPostgresContainer(ctx)
	if err != nil {
		t.Fatalf("Could not start postgres container: %s", err)
	}
	defer cleanup()

	// Prepopulate the database with test data
	_, err = db.Exec(`INSERT INTO posts (title, content, created_at, updated_at) VALUES ('Existing Post', 'Existing content', now(), now())`)
	if err != nil {
		t.Fatalf("Failed to insert initial data: %s", err)
	}

	app := &Application{DB: &database.PostgresDBRepo{DB: db}}

	t.Run("Retrieve Posts Success", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/posts", nil)
		rr := httptest.NewRecorder()

		app.HandleGetPosts(rr, req)

		if status := rr.Code; status != http.StatusOK {
			// Log the response body for debugging
			t.Logf("Response Body: %s", rr.Body.String())
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var posts []*models.Post
		if err := json.NewDecoder(rr.Body).Decode(&posts); err != nil {
			t.Fatalf("Could not decode response: %v", err)
		}

		if len(posts) == 0 {
			t.Errorf("Expected at least one post, got %d", len(posts))
		}
	})
}

func TestHandleUpdatePost_Integration(t *testing.T) {
	ctx := context.Background()
	db, cleanup, err := setupPostgresContainer(ctx)
	if err != nil {
		t.Fatalf("Could not start postgres container: %s", err)
	}
	defer cleanup()

	// Prepopulate the database with test data
	_, err = db.Exec(`INSERT INTO posts (title, content, created_at, updated_at) VALUES ('Existing Post', 'Existing content', now(), now())`)
	if err != nil {
		t.Fatalf("Failed to insert initial data: %s", err)
	}

	app := &Application{DB: &database.PostgresDBRepo{DB: db}}

	t.Run("Update Post Success", func(t *testing.T) {
		body := `{"title":"Updated Title", "content":"Updated content"}`
		req := httptest.NewRequest("PUT", "/posts/1", bytes.NewBufferString(body))
		rr := httptest.NewRecorder()

		// Mocking chi URLParam
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

		app.HandleUpdatePost(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var post models.Post
		if err := json.NewDecoder(rr.Body).Decode(&post); err != nil {
			t.Fatalf("Could not decode response: %v", err)
		}

		if post.Title != "Updated Title" {
			t.Errorf("Handler returned unexpected title: got %v want %v", post.Title, "Updated Title")
		}
	})
}

func TestHandleDeletePost_Integration(t *testing.T) {
	ctx := context.Background()
	db, cleanup, err := setupPostgresContainer(ctx)
	if err != nil {
		t.Fatalf("Could not start postgres container: %s", err)
	}
	defer cleanup()

	// Prepopulate the database with test data
	_, err = db.Exec("INSERT INTO posts (title, content) VALUES ($1, $2)", "To Delete", "Delete me")
	if err != nil {
		t.Fatalf("Failed to insert initial data: %s", err)
	}

	app := &Application{DB: &database.PostgresDBRepo{DB: db}}

	t.Run("Delete Post Success", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/posts/2", nil)
		rr := httptest.NewRecorder()

		// Mocking chi URLParam
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", "2")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

		app.HandleDeletePost(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})
}
