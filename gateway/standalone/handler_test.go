package standalone

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/liquiid727/pipeline-auth/gateway/core"
)

func TestInjectIdentityHeaders_OverridesClientHeaders(t *testing.T) {
	principal := &core.Principal{
		Subject: "alice",
		Owner:   "admin",
		Email:   "alice@example.com",
		Scope:   []string{"read", "write"},
	}
	handler := InjectIdentityHeaders(principal, "trace-1")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("X-Auth-User"); got != "alice" {
			t.Fatalf("expected X-Auth-User alice, got %q", got)
		}
		if got := r.Header.Get("X-Auth-Owner"); got != "admin" {
			t.Fatalf("expected X-Auth-Owner admin, got %q", got)
		}
		if got := r.Header.Get("X-Auth-Email"); got != "alice@example.com" {
			t.Fatalf("expected X-Auth-Email, got %q", got)
		}
		if got := r.Header.Get("X-Auth-Scope"); got != "read write" {
			t.Fatalf("expected scope header, got %q", got)
		}
		if got := r.Header.Get("X-Trace-Id"); got != "trace-1" {
			t.Fatalf("expected trace header, got %q", got)
		}
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/upstream", nil)
	req.Header.Set("X-Auth-User", "mallory")
	req.Header.Set("X-Auth-Owner", "attacker")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected handler to run, got %d", rec.Code)
	}
}
