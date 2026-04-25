package api

import (
	"net/http"
	"github.com/yookibooki/erp-api/internal/app"
	"github.com/yookibooki/erp-api/internal/auth"
)

// TokenVerifier is the ONLY contract the pure API needs.
type TokenVerifier interface {
	Verify(r *http.Request) (auth.Identity, error)
	ResolveRequestContext(r *http.Request) (app.RequestContext, error)
}
