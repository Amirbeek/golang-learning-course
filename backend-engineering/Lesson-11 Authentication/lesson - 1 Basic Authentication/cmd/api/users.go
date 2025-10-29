package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/amirbeek/social/internal/store"
	"github.com/go-chi/chi/v5"
)

type userKey string

const usrCtxKey userKey = "user"

// GetUser godoc
//
//	@Summary		Fetches a user profile
//	@Description	Fetches a user profile by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	store.User
//	@Failure		400	{object}	error	"User payload missing"
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{id} [get]
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user, _ := getUserFromContext(r)

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

type followUser struct {
	UserId int64 `json:"user_id"`
}

// FollowUser godoc
//
//	@Summary		Follows a user
//	@Description	Follows a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int		true	"User ID"
//	@Success		204		{string}	string	"User followed"
//	@Failure		400		{object}	error	"User payload missing"
//	@Failure		404		{object}	error	"User not found"
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/follow [put]
func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser, _ := getUserFromContext(r)

	// TODO: REVERT BACK TO AUTH USERID
	var payload followUser
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	ctx := r.Context()
	err := app.store.Follower.Follow(ctx, followerUser.ID, payload.UserId)
	if err != nil {
		app.internalServeError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusNoContent, nil); err != nil {
		app.internalServeError(w, r, err)
		return
	}
}

// UnfollowUser gdoc
//
//	@Summary		Unfollow a user
//	@Description	Unfollow a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int		true	"User ID"
//	@Success		204		{string}	string	"User unfollowed"
//	@Failure		400		{object}	error	"User payload missing"
//	@Failure		404		{object}	error	"User not found"
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/unfollow [put]
func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	unfollowedUser, _ := getUserFromContext(r)

	// TODO: REVERT BACK TO AUTH USERID
	var payload followUser
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	ctx := r.Context()
	err := app.store.Follower.UnFollow(ctx, unfollowedUser.ID, payload.UserId)
	if err != nil {
		app.internalServeError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, nil); err != nil {
		app.internalServeError(w, r, err)
		return
	}
}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		ctx = context.WithValue(ctx, usrCtxKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromContext(r *http.Request) (*store.User, bool) {
	v := r.Context().Value(usrCtxKey)
	u, _ := v.(*store.User)
	return u, u != nil
}

// ActivateUser godoc
//
//	@Summary		Activates/Register a user
//	@Description	Activates/Register a user by invitation token
//	@Tags			users
//	@Produce		json
//	@Param			token	path		string	true	"Invitation token"
//	@Success		204		{string}	string	"User activated"
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/activate/{token} [put]
func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	err := app.store.Users.Activate(r.Context(), token)

	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundResponse(w, r, err)
		default:
			app.internalServeError(w, r, err)

		}
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServeError(w, r, err)
		return
	}
}
