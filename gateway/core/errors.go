package core

import "fmt"

type AuthErrorCode string

const (
	AuthErrorInvalidState     AuthErrorCode = "invalid_state"
	AuthErrorUnsafeRedirect   AuthErrorCode = "unsafe_redirect"
	AuthErrorMissingSecret    AuthErrorCode = "missing_secret"
	AuthErrorInvalidToken     AuthErrorCode = "invalid_token"
	AuthErrorPermissionDenied AuthErrorCode = "permission_denied"
)

type AuthError struct {
	Code       AuthErrorCode
	StatusCode int
	Message    string
}

func (e *AuthError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func newAuthError(code AuthErrorCode, statusCode int, message string) *AuthError {
	return &AuthError{Code: code, StatusCode: statusCode, Message: message}
}
