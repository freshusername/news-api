package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/freshusername/news-api/models"
	"github.com/freshusername/news-api/validation"
	"github.com/go-chi/chi/v5"
)

func (app *application) HandleCreatePost(w http.ResponseWriter, r *http.Request) {
	post := new(models.Post)

	// Decode the request body into post
	err := json.NewDecoder(r.Body).Decode(post)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	//validate
	validator := validation.NewValidator()
	validator.AddRule("Title", validation.Required())
	validator.AddRule("Title", validation.Length(1, 255))
	validator.AddRule("Content", validation.Required())
	validator.AddRule("Content", validation.Length(1, 500))

	errors := validator.Validate(post)

	if len(errors) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		for _, err := range errors {
			fmt.Fprintf(w, "%s\n", err.PrintError())
		}
		return
	}

	// Close the request body to prevent resource leaks
	defer r.Body.Close()

	createdItem, err := app.DB.CreatePost(post)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusCreated, createdItem)
}

func (app *application) HandleGetPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := app.DB.GetAllPosts()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, posts)
}

func (app *application) HandleUpdatePost(w http.ResponseWriter, r *http.Request) {
	// Extract the item ID from the URL using Chi's URLParam function
	idString := chi.URLParam(r, "id")
	if idString == "" {
		app.errorJSON(w, errors.New("missing item id"), http.StatusBadRequest)
		return
	}

	// Convert the ID from string to int32
	id, err := strconv.ParseInt(idString, 10, 32)
	if err != nil {
		app.errorJSON(w, errors.New("invalid item id format"), http.StatusBadRequest)
		return
	}

	// Decode the request body into a Post struct
	var post *models.Post
	err = json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	//validate
	validator := validation.NewValidator()
	validator.AddRule("Title", validation.Required())
	validator.AddRule("Title", validation.Length(1, 255))
	validator.AddRule("Content", validation.Required())
	validator.AddRule("Content", validation.Length(1, 500))

	errors := validator.Validate(post)

	if len(errors) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		for _, err := range errors {
			fmt.Fprintf(w, "%s\n", err.PrintError())
		}
		return
	}

	// Update the post in the database
	updatedPost, err := app.DB.UpdatePost(int32(id), post)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Prepare a response
	app.writeJSON(w, http.StatusOK, updatedPost)
}

func (app *application) HandleDeletePost(w http.ResponseWriter, r *http.Request) {
	// Extract the item ID from the URL using Chi's URLParam function
	idString := chi.URLParam(r, "id")
	if idString == "" {
		app.errorJSON(w, errors.New("missing item id"), http.StatusBadRequest)
		return
	}

	// Convert the ID from string to int32
	id, err := strconv.ParseInt(idString, 10, 32)
	if err != nil {
		app.errorJSON(w, errors.New("invalid item id format"), http.StatusBadRequest)
		return
	}

	deletedID, err := app.DB.DeletePost(int32(id))
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Prepare a response
	resp := map[string]int32{"id": deletedID}

	app.writeJSON(w, http.StatusOK, resp)
}
