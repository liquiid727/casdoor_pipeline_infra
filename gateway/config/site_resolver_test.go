package config

import "testing"

func TestStaticResolver_ResolvesDomainCaseInsensitive(t *testing.T) {
	resolver := NewStaticResolver([]SiteAuthConfig{
		{
			SiteID:       "admin/site",
			Domain:       "example.com",
			OtherDomains: []string{"alt.example.com"},
			Auth:         AuthConfig{ClientID: "client"},
			Upstream:     UpstreamConfig{URL: "http://upstream.local"},
		},
	})

	site, err := resolver.ResolveByHost("ALT.EXAMPLE.COM:8443")
	if err != nil {
		t.Fatalf("ResolveByHost() error = %v", err)
	}
	if site.SiteID != "admin/site" {
		t.Fatalf("expected site admin/site, got %q", site.SiteID)
	}
}

func TestStaticResolver_RejectsInactiveSite(t *testing.T) {
	resolver := NewStaticResolver([]SiteAuthConfig{
		{SiteID: "admin/site", Domain: "example.com", Status: "Inactive"},
	})

	site, err := resolver.ResolveByHost("example.com")
	if err == nil {
		t.Fatalf("expected inactive site error, got site %#v", site)
	}
}
