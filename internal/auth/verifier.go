package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yookibooki/erp-api/internal/app"
)

type Identity struct {
	UserID      string
	Email       string
	DisplayName string
	OrgID       string
}

type Verifier struct {
	jwksURL string
	issuer  string
}

func NewVerifier(jwksURL, issuer string) (*Verifier, error) {
	return &Verifier{jwksURL: jwksURL, issuer: issuer}, nil
}

func (v *Verifier) Verify(r *http.Request) (Identity, error) {
	authz := r.Header.Get("Authorization")
	if !strings.HasPrefix(strings.ToLower(authz), "bearer ") {
		return Identity{}, errors.New("missing bearer token")
	}
	tokenStr := strings.TrimSpace(authz[7:])

	// TODO: validate against JWKS. For now, parse without verification (DEV ONLY)
	parser := jwt.NewParser()
	token, _, err := parser.ParseUnverified(tokenStr, jwt.MapClaims{})
	if err != nil {
		return Identity{}, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return Identity{}, errors.New("invalid claims")
	}

	return Identity{
		UserID:      fmt.Sprint(claims["sub"]),
		Email:       fmt.Sprint(claims["email"]),
		DisplayName: fmt.Sprint(claims["name"]),
		OrgID:       fmt.Sprint(claims["org_id"]),
	}, nil
}

func (v *Verifier) ResolveRequestContext(r *http.Request) (app.RequestContext, error) {
	id, err := v.Verify(r)
	if err != nil {
		return app.RequestContext{}, err
	}
	orgSlug := r.Header.Get("X-Organization-Slug")
	if orgSlug == "" {
		return app.RequestContext{}, errors.New("X-Organization-Slug required")
	}
	return app.RequestContext{
		UserID:           id.UserID,
		Email:            id.Email,
		DisplayName:      id.DisplayName,
		OrganizationID:   id.OrgID,
		OrganizationSlug: orgSlug,
	}, nil
}
