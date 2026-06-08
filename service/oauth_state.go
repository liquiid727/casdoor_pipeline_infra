package service

import (
	"time"

	"github.com/liquiid727/pipeline-auth/gateway/core"
)

const oauthStateTTL = 10 * time.Minute

type oauthStatePayload struct {
	Path  string `json:"path"`
	Nonce string `json:"nonce"`
	Iat   int64  `json:"iat"`
	Exp   int64  `json:"exp"`
}

func createOAuthState(path string, secret string, now time.Time) (string, error) {
	codec := core.NewStateCodec(secret, core.WithNow(func() time.Time { return now }))
	return codec.Create(path)
}

func parseOAuthState(state string, secret string, now time.Time) (*oauthStatePayload, error) {
	codec := core.NewStateCodec(secret, core.WithNow(func() time.Time { return now }))
	payload, err := codec.Parse(state)
	if err != nil {
		return nil, err
	}
	return &oauthStatePayload{
		Path:  payload.Path,
		Nonce: payload.Nonce,
		Iat:   payload.Iat,
		Exp:   payload.Exp,
	}, nil
}
