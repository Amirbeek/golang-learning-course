package main

import (
	"errors"
	"net/http"
)

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	// pagination, filters, sort

	ctx := r.Context()
	userID := int64(4)
	feed, err := app.store.Posts.GetUserFeed(ctx, userID)

	if err != nil {
		app.internalServeError(w, r, errors.New("error getting user feed"))
		return
	}
	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServeError(w, r, err)
	}
}
