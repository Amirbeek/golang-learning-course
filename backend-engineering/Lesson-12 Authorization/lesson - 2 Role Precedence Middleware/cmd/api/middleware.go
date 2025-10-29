package main

import (
	"context"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func (app *application) AuthTokenMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.unauthorizedErrorResponse(w, r, fmt.Errorf("authorization header missing"))
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			app.unauthorizedErrorResponse(w, r, fmt.Errorf("authorization header format must be 'Bearer <token>'"))
			return
		}

		tokenString := parts[1]
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

		type contextKey string
		const userContextKey = contextKey("user")

		ctx := context.WithValue(r.Context(), userContextKey, user)
		handler.ServeHTTP(w, r.WithContext(ctx))
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
