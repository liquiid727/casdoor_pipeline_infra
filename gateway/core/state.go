package core

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

const defaultStateTTL = 10 * time.Minute

type StateCodec struct {
	secret string
	now    func() time.Time
	ttl    time.Duration
}

type StatePayload struct {
	Path  string `json:"path"`
	Nonce string `json:"nonce"`
	Iat   int64  `json:"iat"`
	Exp   int64  `json:"exp"`
}

func NewStateCodec(secret string, opts ...Option) *StateCodec {
	options := newOptions(opts...)
	return &StateCodec{secret: secret, now: options.now, ttl: options.stateTTL}
}

func (c *StateCodec) Create(path string) (string, error) {
	if c.secret == "" {
		return "", newAuthError(AuthErrorMissingSecret, 400, "state secret is empty")
	}
	if !isSafeRedirectPath(path) {
		return "", newAuthError(AuthErrorUnsafeRedirect, 400, "redirect path must be relative")
	}

	nonceBytes := make([]byte, 16)
	if _, err := rand.Read(nonceBytes); err != nil {
		return "", err
	}

	now := c.now()
	payload := StatePayload{
		Path:  path,
		Nonce: base64.RawURLEncoding.EncodeToString(nonceBytes),
		Iat:   now.Unix(),
		Exp:   now.Add(c.ttl).Unix(),
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	encodedPayload := base64.RawURLEncoding.EncodeToString(payloadBytes)
	return fmt.Sprintf("%s.%s", encodedPayload, c.sign(encodedPayload)), nil
}

func (c *StateCodec) Parse(state string) (*StatePayload, error) {
	if c.secret == "" {
		return nil, newAuthError(AuthErrorMissingSecret, 400, "state secret is empty")
	}

	parts := strings.Split(state, ".")
	if len(parts) != 2 {
		return nil, newAuthError(AuthErrorInvalidState, 400, "state format is invalid")
	}

	expectedSignature := c.sign(parts[0])
	if !hmac.Equal([]byte(expectedSignature), []byte(parts[1])) {
		return nil, newAuthError(AuthErrorInvalidState, 401, "state signature is invalid")
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, err
	}

	payload := &StatePayload{}
	if err := json.Unmarshal(payloadBytes, payload); err != nil {
		return nil, err
	}
	if payload.Exp < c.now().Unix() {
		return nil, newAuthError(AuthErrorInvalidState, 400, "state is expired")
	}
	if !isSafeRedirectPath(payload.Path) {
		return nil, newAuthError(AuthErrorUnsafeRedirect, 400, "redirect path must be relative")
	}
	return payload, nil
}

func (c *StateCodec) sign(encodedPayload string) string {
	mac := hmac.New(sha256.New, []byte(c.secret))
	_, _ = mac.Write([]byte(encodedPayload))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func isSafeRedirectPath(path string) bool {
	return strings.HasPrefix(path, "/") && !strings.HasPrefix(path, "//")
}

func IsAuthError(err error, code AuthErrorCode) bool {
	authErr := &AuthError{}
	return errors.As(err, &authErr) && authErr.Code == code
}
