package standalone

import (
	"net/http"
	"strings"

	"github.com/liquiid727/pipeline-auth/gateway/core"
)

var identityHeaders = []string{
	"X-Auth-User",
	"X-Auth-Owner",
	"X-Auth-Email",
	"X-Auth-Scope",
	"X-Trace-Id",
}

func InjectIdentityHeaders(principal *core.Principal, traceID string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, header := range identityHeaders {
				r.Header.Del(header)
			}
			if principal != nil {
				r.Header.Set("X-Auth-User", principal.Subject)
				r.Header.Set("X-Auth-Owner", principal.Owner)
				r.Header.Set("X-Auth-Email", principal.Email)
				r.Header.Set("X-Auth-Scope", strings.Join(principal.Scope, " "))
			}
			if traceID != "" {
				r.Header.Set("X-Trace-Id", traceID)
			}
			next.ServeHTTP(w, r)
		})
	}
}
