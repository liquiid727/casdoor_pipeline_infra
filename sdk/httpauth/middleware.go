package httpauth

import (
	"errors"
	"net/http"
	"strings"
)

var ErrMissingBearerToken = errors.New("missing bearer token")

func RequireLogin(options Options) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := bearerToken(r)
			if err != nil {
				options.writeUnauthorized(w, r, err)
				return
			}
			if options.Verifier == nil {
				options.writeUnauthorized(w, r, errors.New("token verifier is not configured"))
				return
			}

			principal, err := options.Verifier.VerifyAccessToken(r.Context(), token)
			if err != nil {
				options.writeUnauthorized(w, r, err)
				return
			}
			next.ServeHTTP(w, r.WithContext(ContextWithPrincipal(r.Context(), principal)))
		})
	}
}

func RequireScope(scope string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			principal, ok := CurrentUser(r.Context())
			if !ok || !contains(principal.Scope, scope) {
				Options{}.writeForbidden(w, r, "missing scope")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			principal, ok := CurrentUser(r.Context())
			if !ok || !contains(principal.Roles, role) {
				Options{}.writeForbidden(w, r, "missing role")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func bearerToken(r *http.Request) (string, error) {
	authorization := strings.TrimSpace(r.Header.Get("Authorization"))
	if authorization == "" {
		return "", ErrMissingBearerToken
	}

	scheme, token, ok := strings.Cut(authorization, " ")
	if !ok || !strings.EqualFold(scheme, "Bearer") || strings.TrimSpace(token) == "" {
		return "", ErrMissingBearerToken
	}
	return strings.TrimSpace(token), nil
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
