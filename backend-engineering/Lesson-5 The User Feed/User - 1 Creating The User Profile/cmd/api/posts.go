package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/amirbeek/social/internal/store"
	"github.com/go-chi/chi/v5"
)

type postKey string

const postCtxKey postKey = "post"

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
	post, _ := getPostFromCtx(r)
	comments, err := app.store.Comments.GetByPostId(r.Context(), post.ID)
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

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")

	postID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.store.Posts.DeleteById(ctx, postID)

	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)

		case errors.Is(err, store.NotRowEffectedError):
			app.notFoundResponse(w, r, err)

		case errors.Is(err, store.DeleteError):
			app.notFoundResponse(w, r, err)

		default:
			app.internalServeError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type UpdatePostPayload struct {
	Title   string `json:"title" validate:"omitempty,max=100"`
	Content string `json:"content" validate:"omitempty,max=1000"`
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	post, _ := getPostFromCtx(r)

	var payload UpdatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if payload.Content != "" {
		post.Content = payload.Content
	}

	if payload.Title != "" {
		post.Title = payload.Title
	}

	if err := app.store.Posts.Update(r.Context(), post); err != nil {
		app.internalServeError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusOK, post); err != nil {
		app.internalServeError(w, r, err)
		return
	}
}

func (app *application) postsContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		idStr := chi.URLParam(r, "id")

		if idStr == "" {
			app.badRequestResponse(w, r, errors.New("post ID is missing from URL"))
			return
		}

		postID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			app.badRequestResponse(w, r, errors.New("invalid post ID format"))
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

		ctx = context.WithValue(ctx, postCtxKey, post)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostFromCtx(r *http.Request) (*store.Post, bool) {
	v := r.Context().Value(postCtxKey) // use the SAME key
	p, _ := v.(*store.Post)
	return p, p != nil
}
