package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/amirbeek/social/internal/store"
	"github.com/go-chi/chi/v5"
)

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags" validate:"required,max=100"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	//fmt.Printf("payload: %+v\n", payload)
	payload.Title = strings.TrimSpace(payload.Title)
	payload.Content = strings.TrimSpace(payload.Content)

	if payload.Title == "" || payload.Content == "" {
		app.internalServeError(w, r, errors.New("title or content is empty"))
		return
	}

	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		UserId:  1,
	}

	if err := app.store.Posts.Create(r.Context(), post); err != nil {
		app.internalServeError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		app.internalServeError(w, r, err)
		return
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got")
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")
	postID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.internalServeError(w, r, err)
		return
	}

	post, err := app.store.Posts.GetById(ctx, postID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			app.notFoundResponse(w, r, err)
		} else {
			app.internalServeError(w, r, err)
		}
		return
	}
	comments, err := app.store.Comments.GetByPostId(ctx, postID)
	if err != nil {
		app.internalServeError(w, r, err)
		return
	}

	post.Comments = comments
	if err := writeJSON(w, http.StatusOK, post); err != nil {
		app.internalServeError(w, r, err)
		return
	}
}
