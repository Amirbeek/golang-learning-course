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
	Title   string `json:"title"`
	Content string `json:"content"`
	tags    []string
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload

	if err := readJSON(w, r, &payload); err != nil {
		writeJsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	//fmt.Printf("payload: %+v\n", payload)
	payload.Title = strings.TrimSpace(payload.Title)
	payload.Content = strings.TrimSpace(payload.Content)

	if payload.Title == "" || payload.Content == "" {
		writeJsonError(w, http.StatusBadRequest, "title and content are required")
		return
	}

	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.tags,
		UserId:  1,
	}

	if err := app.store.Posts.Create(r.Context(), post); err != nil {
		writeJsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		writeJsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got")
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")
	postID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJsonError(w, http.StatusBadRequest, "invalid post id")
		return
	}

	post, err := app.store.Posts.GetById(ctx, postID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeJsonError(w, http.StatusNotFound, "post not found")
		} else {
			writeJsonError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	if err := writeJSON(w, http.StatusOK, post); err != nil {
		writeJsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
