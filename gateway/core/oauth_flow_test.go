package core

import (
	"net/url"
	"testing"
	"time"

	gatewayconfig "github.com/liquiid727/pipeline-auth/gateway/config"
)

func TestBeginAuth_BuildsAuthorizeRedirectWithSignedState(t *testing.T) {
	flow := NewOAuthFlow(WithNow(func() time.Time {
		return time.Unix(1_700_000_000, 0)
	}))
	site := gatewayconfig.SiteAuthConfig{
		Auth: gatewayconfig.AuthConfig{
			Endpoint:     "https://auth.example.com",
			ClientID:     "client-id",
			ClientSecret: "client-secret",
		},
	}

	decision, err := flow.BeginAuth(site, "https://app.example.com/caswaf-handler", "/private?x=1")
	if err != nil {
		t.Fatalf("BeginAuth() error = %v", err)
	}
	if decision.Action != AuthActionRedirect {
		t.Fatalf("expected redirect decision, got %q", decision.Action)
	}

	redirectURL, err := url.Parse(decision.RedirectURL)
	if err != nil {
		t.Fatalf("invalid redirect URL: %v", err)
	}
	if redirectURL.Host != "auth.example.com" {
		t.Fatalf("expected auth host, got %q", redirectURL.Host)
	}
	state := redirectURL.Query().Get("state")
	if state == "" || state == "/private?x=1" {
		t.Fatalf("expected signed state, got %q", state)
	}

	session, err := flow.HandleCallback(site, state)
	if err != nil {
		t.Fatalf("HandleCallback() error = %v", err)
	}
	if session.OriginalPath != "/private?x=1" {
		t.Fatalf("expected original path, got %q", session.OriginalPath)
	}
}

func TestBeginAuth_RejectsExternalOriginalPath(t *testing.T) {
	flow := NewOAuthFlow()
	site := gatewayconfig.SiteAuthConfig{Auth: gatewayconfig.AuthConfig{ClientSecret: "secret"}}

	if _, err := flow.BeginAuth(site, "https://app.example.com/caswaf-handler", "https://evil.example"); err == nil {
		t.Fatal("expected external original path to be rejected")
	}
}
