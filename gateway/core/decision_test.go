package core

import "testing"

func TestAuthDecision_Constructors(t *testing.T) {
	principal := &Principal{Subject: "alice", Owner: "admin"}

	allow := Allow(principal)
	if allow.Action != AuthActionAllow || allow.Principal != principal {
		t.Fatalf("unexpected allow decision: %#v", allow)
	}

	redirect := Redirect("https://auth.example.com/login")
	if redirect.Action != AuthActionRedirect || redirect.RedirectURL == "" {
		t.Fatalf("unexpected redirect decision: %#v", redirect)
	}

	deny := Deny(403, "missing scope")
	if deny.Action != AuthActionDeny || deny.StatusCode != 403 || deny.Reason != "missing scope" {
		t.Fatalf("unexpected deny decision: %#v", deny)
	}
}
