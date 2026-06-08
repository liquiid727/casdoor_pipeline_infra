package object

import "testing"

func TestGetSiteByDomain_MatchesCaseInsensitiveDomain(t *testing.T) {
	previous := SiteMap
	defer func() {
		SiteMap = previous
	}()

	site := &Site{Owner: "admin", Name: "site", Domain: "example.com"}
	SiteMap = map[string]*Site{"example.com": site}

	if got := GetSiteByDomain("EXAMPLE.COM"); got != site {
		t.Fatalf("expected case-insensitive domain lookup to return site, got %#v", got)
	}
}

func TestSite_GetChallengeMap(t *testing.T) {
	site := &Site{Challenges: []string{"token:key-auth", "MP_verify_token.txt:verify-body"}}

	got := site.GetChallengeMap()

	if got["token"] != "key-auth" {
		t.Fatalf("expected ACME key auth, got %q", got["token"])
	}
	if got["MP_verify_token.txt"] != "verify-body" {
		t.Fatalf("expected MP verify body, got %q", got["MP_verify_token.txt"])
	}
}

func TestSite_GetHost(t *testing.T) {
	if got := (&Site{Host: "https://upstream.example"}).GetHost(); got != "https://upstream.example" {
		t.Fatalf("expected explicit host, got %q", got)
	}
	if got := (&Site{Port: 7001}).GetHost(); got != "http://localhost:7001" {
		t.Fatalf("expected localhost port host, got %q", got)
	}
}
