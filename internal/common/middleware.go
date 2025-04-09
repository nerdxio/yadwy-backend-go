package common

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

type AuthKey struct{}

func GetAuthMiddlewareFunc(generator *JWTGenerator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := verifyClaimsFromAuthHeader(r, generator)
			if err != nil {
				handleError(w, err)
				return
			}
			ctx := context.WithValue(r.Context(), AuthKey{}, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetAdminMiddlewareFun(generator *JWTGenerator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			claims, err := verifyClaimsFromAuthHeader(r, generator)
			if err != nil {
				handleError(w, err)
				return
			}

			if claims.Role != "ADMIN" {
				handleError(w, NewErrorf(InvalidUserRoleErrorCode, "user is not an admin"))
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
		return nil, NewErrorf(AuthHeaderMissingErrorCode, "authorization header is missing")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		return nil, NewErrorf(AuthHeaderTokenMissingErrorCode, "authorization token is missing")
	}

	claims, err := generator.VerifyToken(token)
	if err != nil {
		return nil, NewErrorf(AuthHeaderTokenVerificationFailed, "authorization token is invalid")
	}

	return claims, nil
}

func GetLoggedInUser(r *http.Request) (*UserClaims, error) {
	claims, ok := r.Context().Value(AuthKey{}).(*UserClaims)
	if !ok {
		return nil, NewErrorf("user is not logged in", "user is not logged in")
	}
	return claims, nil
}

func handleError(w http.ResponseWriter, err error) {
	var appErr *Error
	if errors.As(err, &appErr) {
		switch appErr.Code() {
		case AuthHeaderMissingErrorCode:
			SendError(w, http.StatusUnauthorized, string(appErr.Code()), appErr.Error())
		case AuthHeaderTokenMissingErrorCode:
			SendError(w, http.StatusUnauthorized, string(appErr.Code()), appErr.Error())
		case AuthHeaderTokenVerificationFailed:
			SendError(w, http.StatusUnauthorized, string(appErr.Code()), appErr.Error())
		case InvalidUserRoleErrorCode:
			SendError(w, http.StatusForbidden, string(appErr.Code()), appErr.Error())
		default:
			SendError(w, http.StatusInternalServerError, "internal-server-error", appErr.Error())
		}

	} else {
		SendError(w, http.StatusInternalServerError, "internal-server-error", err.Error())
	}
}
