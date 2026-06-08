package core

import "context"

type TokenVerifier interface {
	VerifyAccessToken(ctx context.Context, token string) (*Principal, error)
}

type TokenVerifierFunc func(ctx context.Context, token string) (*Principal, error)

func (f TokenVerifierFunc) VerifyAccessToken(ctx context.Context, token string) (*Principal, error) {
	return f(ctx, token)
}
