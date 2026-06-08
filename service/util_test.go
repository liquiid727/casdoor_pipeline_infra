package service

import "testing"

func TestJoinPath(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		want string
	}{
		{name: "single slash", a: "https://example.com/", b: "/app", want: "https://example.com/app"},
		{name: "missing slash", a: "https://example.com", b: "app", want: "https://example.com/app"},
		{name: "already joined", a: "https://example.com", b: "/app", want: "https://example.com/app"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := joinPath(tt.a, tt.b); got != tt.want {
				t.Fatalf("joinPath(%q, %q) = %q, want %q", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestGetHostNonWww(t *testing.T) {
	if got := getHostNonWww("www.example.com"); got != "example.com" {
		t.Fatalf("expected example.com, got %q", got)
	}
	if got := getHostNonWww("example.com"); got != "" {
		t.Fatalf("expected empty host for non-www input, got %q", got)
	}
}

func TestGetDomainWithoutPort(t *testing.T) {
	if got := getDomainWithoutPort("example.com:8443"); got != "example.com" {
		t.Fatalf("expected example.com, got %q", got)
	}
	if got := getDomainWithoutPort("example.com"); got != "example.com" {
		t.Fatalf("expected example.com, got %q", got)
	}
}
