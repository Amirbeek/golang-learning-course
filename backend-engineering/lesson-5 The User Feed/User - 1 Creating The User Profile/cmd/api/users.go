package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/amirbeek/social/internal/store"
	"github.com/go-chi/chi/v5"
)

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")

	postID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user, err := app.store.Users.GetUserById(ctx, postID)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.badRequestResponse(w, r, err)
		}
	}

	if err := writeJSON(w, http.StatusOK, user); err != nil {
		app.internalServeError(w, r, err)
		return
	}

}

func (app *application) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")
	userId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.store.Users.DeleteById(ctx, userId)

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
