package common

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

type AuthKey struct{}

func GetAuthMiddlewareFunc(generator *JWTGenerator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := verifyClaimsFromAuthHeader(r, generator)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), AuthKey{}, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func verifyClaimsFromAuthHeader(r *http.Request, generator *JWTGenerator) (*UserClaims, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("authorization header is missing")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		return nil, fmt.Errorf("bearer token is missing")
	}

	claims, err := generator.VerifyToken(token)
	if err != nil {
		return nil, fmt.Errorf("error verifying token: %w", err)
	}

	return claims, nil
}

func GetLoggedInUser(r *http.Request) (*UserClaims, error) {
	claims, ok := r.Context().Value(AuthKey{}).(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("user claims not found in context")
	}
	return claims, nil
}
