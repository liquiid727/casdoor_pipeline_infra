package service

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetScheme_UsesForwardedProto(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://example.com/app", nil)
	req.Header.Set("X-Forwarded-Proto", "https")

	if got := getScheme(req); got != "https" {
		t.Fatalf("expected forwarded proto https, got %q", got)
	}
}

func TestGetScheme_UsesTLS(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "https://example.com/app", nil)

	if got := getScheme(req); got != "https" {
		t.Fatalf("expected TLS request scheme https, got %q", got)
	}
}

func TestGetScheme_UsesForwardedScheme(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://example.com/app", nil)
	req.Header.Set("X-Forwarded-Scheme", "https")

	if got := getScheme(req); got != "https" {
		t.Fatalf("expected forwarded scheme https, got %q", got)
	}
}

func TestGetScheme_UsesFirstForwardedValue(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://example.com/app", nil)
	req.Header.Set("X-Forwarded-Proto", "https, http")

	if got := getScheme(req); got != "https" {
		t.Fatalf("expected first forwarded value https, got %q", got)
	}
}

func TestOAuthState_RoundTrip(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	state, err := createOAuthState("/private?x=1", "secret", now)
	if err != nil {
		t.Fatalf("createOAuthState() error = %v", err)
	}

	parsed, err := parseOAuthState(state, "secret", now.Add(time.Minute))
	if err != nil {
		t.Fatalf("parseOAuthState() error = %v", err)
	}
	if parsed.Path != "/private?x=1" {
		t.Fatalf("expected path to round trip, got %q", parsed.Path)
	}
	if parsed.Nonce == "" {
		t.Fatal("expected nonce to be populated")
	}
}

func TestOAuthState_RejectsTamperedPayload(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	state, err := createOAuthState("/private", "secret", now)
	if err != nil {
		t.Fatalf("createOAuthState() error = %v", err)
	}

	tampered := state[:len(state)-1] + "x"
	if _, err := parseOAuthState(tampered, "secret", now.Add(time.Minute)); err == nil {
		t.Fatal("expected tampered state to be rejected")
	}
}

func TestOAuthState_RejectsExpiredPayload(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	state, err := createOAuthState("/private", "secret", now)
	if err != nil {
		t.Fatalf("createOAuthState() error = %v", err)
	}

	if _, err := parseOAuthState(state, "secret", now.Add(11*time.Minute)); err == nil {
		t.Fatal("expected expired state to be rejected")
	}
}

func TestOAuthState_RejectsExternalRedirect(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)

	if _, err := createOAuthState("https://evil.example/path", "secret", now); err == nil {
		t.Fatal("expected absolute URL redirect to be rejected")
	}
	if _, err := createOAuthState("//evil.example/path", "secret", now); err == nil {
		t.Fatal("expected protocol-relative redirect to be rejected")
	}
}

func TestOAuthState_RejectsEmptySecret(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)

	if _, err := createOAuthState("/private", "", now); err == nil {
		t.Fatal("expected empty secret to be rejected when creating state")
	}
	state, err := createOAuthState("/private", "secret", now)
	if err != nil {
		t.Fatalf("createOAuthState() error = %v", err)
	}
	if _, err := parseOAuthState(state, "", now); err == nil {
		t.Fatal("expected empty secret to be rejected when parsing state")
	}
}
