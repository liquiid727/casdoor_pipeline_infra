package service

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/liquiid727/pipeline-auth/object"
)

func withSiteMap(t *testing.T, siteMap map[string]*object.Site) {
	t.Helper()

	previous := object.SiteMap
	object.SiteMap = siteMap
	t.Cleanup(func() {
		object.SiteMap = previous
	})
}

func newGatewayRequest(uri string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, uri, nil)
	req.Host = "example.com"
	req.Header.Set("User-Agent", "Uptime-Kuma")
	return req
}

func TestHandleRequest_InactiveSiteReturns503(t *testing.T) {
	withSiteMap(t, map[string]*object.Site{
		"example.com": {
			Owner:  "admin",
			Name:   "site",
			Domain: "example.com",
			Host:   "http://upstream.local",
			Status: "Inactive",
		},
	})

	rec := httptest.NewRecorder()
	handleRequest(rec, newGatewayRequest("/private"))

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503 for inactive site, got %d body %q", rec.Code, rec.Body.String())
	}
}

func TestHandleRequest_AcmeChallengeStillWorks(t *testing.T) {
	withSiteMap(t, map[string]*object.Site{
		"example.com": {
			Owner:      "admin",
			Name:       "site",
			Domain:     "example.com",
			Status:     "Inactive",
			Challenges: []string{"token:key-auth"},
		},
	})

	rec := httptest.NewRecorder()
	handleRequest(rec, newGatewayRequest("/.well-known/acme-challenge/token"))

	if rec.Code != http.StatusOK {
		t.Fatalf("expected ACME challenge to bypass inactive block, got %d", rec.Code)
	}
	if rec.Body.String() != "key-auth" {
		t.Fatalf("expected key authorization body, got %q", rec.Body.String())
	}
}

func TestHandleRequest_MpVerifyStillWorks(t *testing.T) {
	withSiteMap(t, map[string]*object.Site{
		"example.com": {
			Owner:      "admin",
			Name:       "site",
			Domain:     "example.com",
			Status:     "Inactive",
			Challenges: []string{"MP_verify_token.txt:verify-body"},
		},
	})

	rec := httptest.NewRecorder()
	handleRequest(rec, newGatewayRequest("/MP_verify_token.txt"))

	if rec.Code != http.StatusOK {
		t.Fatalf("expected MP verify to bypass inactive block, got %d", rec.Code)
	}
	if rec.Body.String() != "verify-body" {
		t.Fatalf("expected verify body, got %q", rec.Body.String())
	}
}

func TestHandleRequest_WwwRedirect(t *testing.T) {
	withSiteMap(t, map[string]*object.Site{
		"example.com": {
			Owner:  "admin",
			Name:   "site",
			Domain: "example.com",
			Host:   "http://upstream.local",
		},
	})
	req := httptest.NewRequest(http.MethodGet, "/private", nil)
	req.Host = "www.example.com"
	req.Header.Set("User-Agent", "Uptime-Kuma")
	rec := httptest.NewRecorder()

	handleRequest(rec, req)

	if rec.Code != http.StatusMovedPermanently {
		t.Fatalf("expected www redirect, got %d", rec.Code)
	}
	if got := rec.Header().Get("Location"); got != "http://example.com/private" {
		t.Fatalf("expected redirect to non-www host, got %q", got)
	}
}

func TestHandleRequest_HTTPSOnlyRedirect(t *testing.T) {
	withSiteMap(t, map[string]*object.Site{
		"example.com": {
			Owner:   "admin",
			Name:    "site",
			Domain:  "example.com",
			Host:    "http://upstream.local",
			SslMode: "HTTPS Only",
		},
	})

	rec := httptest.NewRecorder()
	handleRequest(rec, newGatewayRequest("/private"))

	if rec.Code != http.StatusMovedPermanently {
		t.Fatalf("expected HTTPS redirect, got %d", rec.Code)
	}
	if got := rec.Header().Get("Location"); got != "https://example.com/private" {
		t.Fatalf("expected redirect to https URL, got %q", got)
	}
}
