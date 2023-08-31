package middleware

import (
	"context"
	"errors"
	"net/http"
	"social-network/pkg/models"
	utils "social-network/pkg/utils"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type contextKey struct {
	id string
}

var UserIDKey = &contextKey{id: "authenticatedUser"}

func AuthenticatedUser(db *sqlx.DB) mux.MiddlewareFunc {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Unable to parse form data", http.StatusBadRequest)
				return
			}

			cookie, err := r.Cookie("accessToken")
			if err != nil {
				if err == http.ErrNoCookie {
					http.Error(w, "Unauthorized access", http.StatusUnauthorized)
					return
				}
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			// Validate claims from jwt cookie
			tkn := strings.TrimPrefix(cookie.Value, "=")
			claims, err := utils.GetClaimsFromJWT(tkn)
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			// Make sure that claims are not malformed and userid and uuid matches in database
			id, err := VerifySession(db, claims)
			if err != nil {
				http.Error(w, "Unauthorized access", http.StatusUnauthorized)
				return
			}

			// // Everything OK, call the next handler with userid. Will be used for websocket connection.
			ctx := context.WithValue(r.Context(), UserIDKey, id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func VerifySession(db *sqlx.DB, claims jwt.MapClaims) (int, error) {
	id, ok := claims["userID"].(float64)
	if !ok {
		return 0, errors.New("userID not found in claims or is of wrong type")
	}

	uuid, ok := claims["uuid"].(string)
	if !ok {
		return 0, errors.New("uuid not found in claims or is of wrong type")
	}

	issuer, ok := claims["iss"].(string)
	if !ok {
		return 0, errors.New("issuer not found in claims or is of wrong type")
	}

	if issuer != "fakebook" {
		return 0, errors.New("malformed jwt claims")
	}

	valid, err := models.ValidateSession(db, int(id), uuid)
	if err != nil || !valid {
		return 0, errors.New("malformed jwt claims")
	}

	return int(id), nil
}
