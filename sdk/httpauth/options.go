package httpauth

import (
	"net/http"

	"github.com/liquiid727/pipeline-auth/gateway/core"
)

type Options struct {
	Verifier             core.TokenVerifier
	UnauthorizedResponse func(http.ResponseWriter, *http.Request, error)
	ForbiddenResponse    func(http.ResponseWriter, *http.Request, string)
}

func (o Options) writeUnauthorized(w http.ResponseWriter, r *http.Request, err error) {
	if o.UnauthorizedResponse != nil {
		o.UnauthorizedResponse(w, r, err)
		return
	}
	http.Error(w, "unauthorized", http.StatusUnauthorized)
}

func (o Options) writeForbidden(w http.ResponseWriter, r *http.Request, reason string) {
	if o.ForbiddenResponse != nil {
		o.ForbiddenResponse(w, r, reason)
		return
	}
	http.Error(w, "forbidden", http.StatusForbidden)
}
