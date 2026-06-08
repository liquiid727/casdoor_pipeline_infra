package core

import (
	"fmt"
	"net/url"
	"time"

	gatewayconfig "github.com/liquiid727/pipeline-auth/gateway/config"
)

type Option func(*options)

type options struct {
	now      func() time.Time
	stateTTL time.Duration
}

func WithNow(now func() time.Time) Option {
	return func(opts *options) {
		opts.now = now
	}
}

func newOptions(opts ...Option) options {
	options := options{
		now:      time.Now,
		stateTTL: defaultStateTTL,
	}
	for _, opt := range opts {
		opt(&options)
	}
	return options
}

type OAuthFlow struct {
	opts []Option
}

func NewOAuthFlow(opts ...Option) *OAuthFlow {
	return &OAuthFlow{opts: opts}
}

func (f *OAuthFlow) BeginAuth(site gatewayconfig.SiteAuthConfig, callbackURL string, originalPath string) (AuthDecision, error) {
	codec := NewStateCodec(site.Auth.ClientSecret, f.opts...)
	state, err := codec.Create(originalPath)
	if err != nil {
		return AuthDecision{}, err
	}

	redirectURL := fmt.Sprintf("%s/login/oauth/authorize?client_id=%s&response_type=code&redirect_uri=%s&scope=read&state=%s",
		site.Auth.Endpoint,
		url.QueryEscape(site.Auth.ClientID),
		url.QueryEscape(callbackURL),
		url.QueryEscape(state),
	)
	return Redirect(redirectURL), nil
}

func (f *OAuthFlow) HandleCallback(site gatewayconfig.SiteAuthConfig, state string) (*AuthSession, error) {
	codec := NewStateCodec(site.Auth.ClientSecret, f.opts...)
	payload, err := codec.Parse(state)
	if err != nil {
		return nil, err
	}
	return &AuthSession{OriginalPath: payload.Path}, nil
}
