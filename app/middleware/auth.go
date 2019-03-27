package middleware

import (
	"context"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"payments/app/models"
	u "payments/utils"
	"strings"
)

const ERROR_MISSING_TOKEN = "Missing auth token"
const ERROR_MALFORMED_TOKEN = "Invalid/Malformed auth token"
const ERROR_TOKEN_INVALID = "Token Invalid"

var JwtAuthentication = func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// List of endpoints that doesn't require auth
		notAuth := []string{"/v1/user", "/v1/user/login", "/v1/health"}

		// Current Request Path
		requestPath := r.URL.Path

		// Check if request does not need authentication, serve the request if it doesn't need it
		for _, value := range notAuth {

			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		// Grab the token from the header
		tokenHeader := r.Header.Get("Authorization")

		// Token is missing, returns with error code 403 Unauthorized
		if tokenHeader == "" {

			if response, err := json.Marshal(u.Response{Errors: []string{ERROR_MISSING_TOKEN}}); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				w.Write(response)
			}
			return
		}

		// The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
		splitted := strings.Split(tokenHeader, " ")
		if len(splitted) != 2 {
			if response, err := json.Marshal(u.Response{Errors: []string{ERROR_MALFORMED_TOKEN}}); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				w.Write(response)
			}
			return
		}

		// Grab the token part, what we are truly interested in
		tokenPart := splitted[1]
		tk := &models.Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})

		// Malformed token, returns with http code 403
		if err != nil {
			if response, err := json.Marshal(u.Response{Errors: []string{ERROR_MALFORMED_TOKEN}}); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				w.Write(response)
			}
			return
		}

		// Token is invalid, maybe not signed on this server
		if !token.Valid {
			if response, err := json.Marshal(u.Response{Errors: []string{ERROR_TOKEN_INVALID}}); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				w.Write(response)
			}
			return
		}

		// Everything is OK, proceed with the request and set the caller to the user retrieved from the parsed token
		ctx := context.WithValue(r.Context(), "user", tk.UserId)
		r = r.WithContext(ctx)

		// Proceed in the middleware chain!
		next.ServeHTTP(w, r)
	})
}
