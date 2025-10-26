package store

import (
	"errors"
	"net/http"

	"github.com/amirbeek/social/cmd/api"
)

func (app *main.application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	// pagination, filters

	fq := main.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "asc",
	}
	fq, err := fq.Parse(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if err := main.Validate.Struct(fq); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
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
