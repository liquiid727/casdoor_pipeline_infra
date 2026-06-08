package httpauth

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/liquiid727/pipeline-auth/gateway/core"
)

func TestRequireLogin_AllowsValidBearerToken(t *testing.T) {
	middleware := RequireLogin(Options{
		Verifier: core.TokenVerifierFunc(func(ctx context.Context, token string) (*core.Principal, error) {
			if token != "valid-token" {
				t.Fatalf("expected valid-token, got %q", token)
			}
			return &core.Principal{Subject: "alice", Owner: "admin", Scope: []string{"read"}, Roles: []string{"user"}}, nil
		}),
	})

	req := httptest.NewRequest(http.MethodGet, "/api", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	rec := httptest.NewRecorder()

	middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		principal, ok := CurrentUser(r.Context())
		if !ok {
			t.Fatal("expected principal in request context")
		}
		if principal.Subject != "alice" {
			t.Fatalf("expected alice principal, got %#v", principal)
		}
		w.WriteHeader(http.StatusNoContent)
	})).ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected handler to run, got %d", rec.Code)
	}
}

func TestRequireLogin_RejectsMissingBearerToken(t *testing.T) {
	middleware := RequireLogin(Options{
		Verifier: core.TokenVerifierFunc(func(ctx context.Context, token string) (*core.Principal, error) {
			t.Fatal("verifier should not be called without bearer token")
			return nil, nil
		}),
	})

	rec := httptest.NewRecorder()
	middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not run")
	})).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api", nil))

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestRequireLogin_RejectsInvalidToken(t *testing.T) {
	middleware := RequireLogin(Options{
		Verifier: core.TokenVerifierFunc(func(ctx context.Context, token string) (*core.Principal, error) {
			return nil, errors.New("invalid token")
		}),
	})
	req := httptest.NewRequest(http.MethodGet, "/api", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	rec := httptest.NewRecorder()

	middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not run")
	})).ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestRequireScope_DeniesMissingScope(t *testing.T) {
	handler := WithPrincipal(&core.Principal{Subject: "alice", Scope: []string{"read"}})(
		RequireScope("write")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("handler should not run")
		})),
	)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api", nil))

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", rec.Code)
	}
}
