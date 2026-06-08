package service

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/liquiid727/pipeline-auth/object"
)

func TestSetAuthCookie_HasSecurityAttributes(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "https://example.com/", nil)
	rec := httptest.NewRecorder()

	setAuthCookie(rec, req, "token")

	cookies := rec.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected one cookie, got %d", len(cookies))
	}

	cookie := cookies[0]
	if cookie.Name != authCookieName {
		t.Fatalf("expected auth cookie name %q, got %q", authCookieName, cookie.Name)
	}
	if cookie.Value != "token" {
		t.Fatalf("expected token value, got %q", cookie.Value)
	}
	if cookie.Path != "/" {
		t.Fatalf("expected path /, got %q", cookie.Path)
	}
	if !cookie.HttpOnly {
		t.Fatal("expected HttpOnly cookie")
	}
	if !cookie.Secure {
		t.Fatal("expected Secure cookie for https request")
	}
	if cookie.SameSite != http.SameSiteLaxMode {
		t.Fatalf("expected SameSite=Lax, got %v", cookie.SameSite)
	}
}

func TestClearAuthCookie_ExpiresCookie(t *testing.T) {
	rec := httptest.NewRecorder()

	clearAuthCookie(rec)

	cookies := rec.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected one cookie, got %d", len(cookies))
	}
	cookie := cookies[0]
	if cookie.Name != authCookieName {
		t.Fatalf("expected auth cookie name %q, got %q", authCookieName, cookie.Name)
	}
	if cookie.MaxAge != -1 {
		t.Fatalf("expected MaxAge=-1, got %d", cookie.MaxAge)
	}
	if cookie.Path != "/" {
		t.Fatalf("expected path /, got %q", cookie.Path)
	}
	if !cookie.HttpOnly {
		t.Fatal("expected HttpOnly cookie")
	}
}

func TestHandleRequest_InvalidTokenClearsCookieAndRedirects(t *testing.T) {
	withSiteMap(t, map[string]*object.Site{
		"example.com": {
			Owner:           "admin",
			Name:            "site",
			Domain:          "example.com",
			Host:            "http://upstream.local",
			AuthApplication: "app",
			ApplicationObj: &object.Application{
				Owner:        "admin",
				Name:         "app",
				ClientId:     "client-id",
				ClientSecret: "client-secret",
			},
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/private", nil)
	req.Host = "example.com"
	req.Header.Set("User-Agent", "Uptime-Kuma")
	req.AddCookie(&http.Cookie{Name: authCookieName, Value: "invalid-token"})
	rec := httptest.NewRecorder()

	handleRequest(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("expected redirect, got status %d body %q", rec.Code, rec.Body.String())
	}
	if got := rec.Header().Get("Location"); got == "" {
		t.Fatal("expected redirect location")
	}

	foundClearedCookie := false
	for _, cookie := range rec.Result().Cookies() {
		if cookie.Name == authCookieName && cookie.MaxAge == -1 {
			foundClearedCookie = true
		}
	}
	if !foundClearedCookie {
		t.Fatalf("expected auth cookie to be cleared, got %v", rec.Result().Cookies())
	}
}
