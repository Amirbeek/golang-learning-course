package main

import (
	"context"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/amirbeek/social/internal/store"
	"github.com/golang-jwt/jwt/v5"
)

func (app *application) checkPostOwnership(role string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, _ := getUserFromContext(r)
		post, _ := getPostFromCtx(r)

		if post.UserId == user.ID {
			next.ServeHTTP(w, r)
			return
		}

		allowed, err := app.checkRolePrecedence(r.Context(), user, role)
		if err != nil {
			app.internalServeError(w, r, err)
			return
		}
		if !allowed {
			app.forbiddenResponse(w, r, nil)
			return
		}

		next.ServeHTTP(w, r)
	}
}
func (app *application) checkRolePrecedence(ctx context.Context, user *store.User, role string) (bool, error) {
	targetRole, err := app.store.Roles.GetByName(ctx, role)

	if err != nil {
		return false, err
	}

	return user.Role.Level >= targetRole.Level, nil
}

func (app *application) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.unauthorizedErrorResponse(w, r, fmt.Errorf("authorization header missing"))
			return
		}
		fmt.Println("FIRST 1: ", authHeader)

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			app.unauthorizedErrorResponse(w, r, fmt.Errorf("authorization header is malformed"))
			return
		}
		fmt.Println("FIRST 2: ", parts)

		// Tolerate accidental quotes if the client copied a JSON string value
		tokenString := parts[1]
		tokenString = strings.Trim(tokenString, "\"")

		fmt.Println("TOKEN: ", tokenString)

		jwtToken, err := app.authenticator.ValidateToken(tokenString)
		if err != nil {
			app.unauthorizedErrorResponse(w, r, err)
			return
		}

		if !jwtToken.Valid {
			app.unauthorizedErrorResponse(w, r, fmt.Errorf("invalid token"))
			return
		}

		claims, ok := jwtToken.Claims.(jwt.MapClaims)
		if !ok {
			app.unauthorizedErrorResponse(w, r, fmt.Errorf("invalid JWT claims"))
			return
		}

		userID, err := strconv.ParseInt(fmt.Sprintf("%v", claims["sub"]), 10, 64)
		if err != nil {
			app.unauthorizedErrorResponse(w, r, fmt.Errorf("invalid user ID in token"))
			return
		}

		user, err := app.store.Users.GetUserById(r.Context(), userID)
		if err != nil {
			app.unauthorizedErrorResponse(w, r, err)
			return
		}

		ctx := context.WithValue(r.Context(), usrCtxKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) BasicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Check the header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.unauthorizedBasicErrorResponse(w, r, fmt.Errorf("authorization header missing"))
			return
		}

		// Validate the format of the Authorization header
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Basic") {
			app.unauthorizedBasicErrorResponse(w, r, fmt.Errorf("authorization header malformed"))
			return
		}

		// Decode the Base64 credentials
		decoded, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			app.unauthorizedBasicErrorResponse(w, r, fmt.Errorf("invalid base64 encoding"))
			return
		}

		// Extract username and password
		creds := strings.SplitN(string(decoded), ":", 2)
		if len(creds) != 2 {
			app.unauthorizedBasicErrorResponse(w, r, fmt.Errorf("invalid credentials format"))
			return
		}

		// Retrieve expected credentials from configuration
		username := app.config.auth.basic.user
		password := app.config.auth.basic.pass

		// Compare credentials securely
		if subtle.ConstantTimeCompare([]byte(creds[0]), []byte(username)) != 1 ||
			subtle.ConstantTimeCompare([]byte(creds[1]), []byte(password)) != 1 {
			app.unauthorizedBasicErrorResponse(w, r, fmt.Errorf("invalid credentials"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
