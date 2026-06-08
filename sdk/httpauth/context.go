package httpauth

import (
	"context"
	"net/http"

	"github.com/liquiid727/pipeline-auth/gateway/core"
)

type principalContextKey struct{}

func ContextWithPrincipal(ctx context.Context, principal *core.Principal) context.Context {
	return context.WithValue(ctx, principalContextKey{}, principal)
}

func CurrentUser(ctx context.Context) (*core.Principal, bool) {
	principal, ok := ctx.Value(principalContextKey{}).(*core.Principal)
	return principal, ok && principal != nil
}

func WithPrincipal(principal *core.Principal) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(ContextWithPrincipal(r.Context(), principal)))
		})
	}
}
