package core

type AuthAction string

const (
	AuthActionAllow    AuthAction = "allow"
	AuthActionRedirect AuthAction = "redirect"
	AuthActionDeny     AuthAction = "deny"
	AuthActionClear    AuthAction = "clear"
)

type AuthDecision struct {
	Action       AuthAction
	StatusCode   int
	Reason       string
	RedirectURL  string
	Principal    *Principal
	ClearSession bool
}

func Allow(principal *Principal) AuthDecision {
	return AuthDecision{Action: AuthActionAllow, Principal: principal}
}

func Redirect(redirectURL string) AuthDecision {
	return AuthDecision{Action: AuthActionRedirect, StatusCode: 302, RedirectURL: redirectURL}
}

func Deny(statusCode int, reason string) AuthDecision {
	return AuthDecision{Action: AuthActionDeny, StatusCode: statusCode, Reason: reason}
}

func ClearAndRedirect(redirectURL string) AuthDecision {
	return AuthDecision{Action: AuthActionRedirect, StatusCode: 302, RedirectURL: redirectURL, ClearSession: true}
}
