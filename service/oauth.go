// Copyright 2023 The casbin Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	gatewayconfig "github.com/liquiid727/pipeline-auth/gateway/config"
	"github.com/liquiid727/pipeline-auth/gateway/core"
	"github.com/liquiid727/pipeline-auth/object"
	"github.com/liquiid727/pipeline-auth/util"
)

func getSignInURL(authServerClient *casdoorsdk.Client, callbackURL string, originalPath string) string {
	scope := "read"
	return fmt.Sprintf("%s/login/oauth/authorize?client_id=%s&response_type=code&redirect_uri=%s&scope=%s&state=%s",
		authServerClient.Endpoint, authServerClient.ClientId, url.QueryEscape(callbackURL), scope, url.QueryEscape(originalPath))
}

func redirectToAuthServer(authServerClient *casdoorsdk.Client, w http.ResponseWriter, r *http.Request) {
	scheme := getScheme(r)

	callbackURL := fmt.Sprintf("%s://%s/caswaf-handler", scheme, r.Host)
	flow := core.NewOAuthFlow()
	decision, err := flow.BeginAuth(gatewayconfig.SiteAuthConfig{
		Auth: gatewayconfig.AuthConfig{
			Endpoint:     authServerClient.Endpoint,
			ClientID:     authServerClient.ClientId,
			ClientSecret: authServerClient.ClientSecret,
		},
	}, callbackURL, getRequestReturnPath(r))
	if err != nil {
		responseError(w, "Pipeline Auth WAF error: create OAuth state failed: %s", err.Error())
		return
	}
	http.Redirect(w, r, decision.RedirectURL, http.StatusFound)
}

func getRequestReturnPath(r *http.Request) string {
	path := r.URL.RequestURI()
	if path != "" {
		return path
	}
	return r.RequestURI
}

func handleAuthCallback(w http.ResponseWriter, r *http.Request) {
	site := getSiteByDomainWithWww(r.Host)
	if site == nil {
		observeGatewayCallbackFailure(nil, "site_not_found")
		responseError(w, "CasWAF error: site not found for host: %s", r.Host)
		return
	}

	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	if code == "" {
		observeGatewayCallbackFailure(site, "missing_code")
		responseError(w, "CasWAF error: the code should not be empty")
		return
	} else if state == "" {
		observeGatewayCallbackFailure(site, "missing_state")
		responseError(w, "CasWAF error: the state should not be empty")
		return
	}

	application, err := object.GetApplication(util.GetId(site.Owner, site.AuthApplication))
	if err != nil {
		observeGatewayCallbackFailure(site, "application_error")
		responseError(w, "Pipeline Auth WAF error: auth server token exchange failed: %s", err.Error())
		return
	}
	if application == nil {
		observeGatewayCallbackFailure(site, "application_not_found")
		http.Error(w, "Pipeline Auth WAF error: auth application not found", http.StatusBadRequest)
		return
	}

	parsedState, err := parseOAuthState(state, application.ClientSecret, time.Now())
	if err != nil {
		observeGatewayCallbackFailure(site, "invalid_state")
		http.Error(w, "Pipeline Auth WAF error: invalid OAuth state", http.StatusBadRequest)
		return
	}

	// authServerClient, err := getAuthServerClientFromSite(site)
	//if err != nil {
	//	responseError(w, "Pipeline Auth WAF error: getAuthServerClientFromSite() error: %s", err.Error())
	//	return
	//}

	token, tokenError, err := object.GetAuthorizationCodeToken(application, application.ClientSecret, code, "", "")
	if tokenError != nil {
		observeGatewayCallbackFailure(site, "token_error")
		responseError(w, "Pipeline Auth WAF error: auth server token exchange failed: %s", tokenError.Error)
		return
	}
	if err != nil {
		observeGatewayCallbackFailure(site, "token_exchange_error")
		responseError(w, "Pipeline Auth WAF error: auth server token exchange failed: %s", err.Error())
		return
	}

	setAuthCookie(w, r, token.AccessToken)
	observeGatewayCallbackSuccess(site)

	http.Redirect(w, r, parsedState.Path, http.StatusFound)
}
