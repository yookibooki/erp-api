package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/yookibooki/erp-api/internal/app"
	"github.com/yookibooki/erp-api/internal/auth"
	"github.com/yookibooki/erp-api/internal/organizations"
)

// TokenVerifier is the ONLY auth contract the pure API needs.
// It verifies a Bearer token and returns an identity.
// No cookies, no redirects, no login handlers - pure stateless API.
type TokenVerifier interface {
	Verify(r *http.Request) (auth.Identity, error)
	ResolveRequestContext(r *http.Request) (app.RequestContext, error)
}

type Dependencies struct {
	Auth                 TokenVerifier
	OrganizationsService *organizations.Service
}

type Server struct {
	auth          TokenVerifier
	organizations *organizations.Service
}

func NewServer(deps Dependencies) *Server {
	return &Server{
		auth:          deps.Auth,
		organizations: deps.OrganizationsService,
	}
}

// Router returns a pure JSON API router - no UI routes, no auth flows.
func (s *Server) Router() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)

	// Health - no auth
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// API v1 - all routes require Bearer token
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(s.requireAuth)

		r.Get("/me", s.handleMe)
		r.Get("/organizations", s.handleOrganizationsList)
		r.Post("/organizations", s.handleOrganizationsCreate)
		r.Get("/organizations/current", s.handleCurrentOrganization)
		// Note: switching org is now client-side via header, not server cookie
	})

	return r
}

func (s *Server) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		identity, err := s.auth.Verify(r)
		if err != nil {
			writeAPIError(w, http.StatusUnauthorized, "unauthenticated", "missing or invalid bearer token", nil)
			return
		}
		// Attach identity to context if needed downstream
		ctx := r.Context()
		next.ServeHTTP(w, r.WithContext(ctx))
		_ = identity
	})
}

func (s *Server) handleMe(w http.ResponseWriter, r *http.Request) {
	identity, err := s.auth.Verify(r)
	if err != nil {
		writeAPIError(w, http.StatusUnauthorized, "unauthenticated", err.Error(), nil)
		return
	}

	rc, err := s.auth.ResolveRequestContext(r)
	if err != nil {
		writeAPIError(w, http.StatusForbidden, "no_organization", err.Error(), nil)
		return
	}

	// Build response - pure data, no CSRF, no cookies
	resp := mapMe(rc, identity.Memberships)
	writeJSON(w, http.StatusOK, resp)
}

func (s *Server) handleOrganizationsList(w http.ResponseWriter, r *http.Request) {
	identity, ok := s.requireIdentity(w, r)
	if !ok {
		return
	}

	items := make([]OrganizationDTO, 0, len(identity.Memberships))
	for _, m := range identity.Memberships {
		items = append(items, mapOrganizationFromMembership(m, false))
	}

	writeJSON(w, http.StatusOK, OrganizationListResponse{
		Items: items,
		Meta:  ListMeta{Total: len(items), Limit: len(items), Offset: 0},
	})
}

func (s *Server) handleOrganizationsCreate(w http.ResponseWriter, r *http.Request) {
	identity, ok := s.requireIdentity(w, r)
	if !ok {
		return
	}

	req, err := parseJSON[CreateOrganizationRequest](r)
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, "invalid_json", err.Error(), nil)
		return
	}

	tenant, err := s.organizations.Create(r.Context(), organizations.CreateInput{
		OwnerUserID:  identity.User.ID,
		CompanyName:  req.Name,
		Slug:         req.Slug,
		Country:      req.Country,
		BaseCurrency: req.BaseCurrency,
	})
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, "create_failed", err.Error(), nil)
		return
	}

	writeJSON(w, http.StatusCreated, CreateOrganizationResponse{
		Organization: mapOrganization(tenant, "owner", true),
	})
}

func (s *Server) handleCurrentOrganization(w http.ResponseWriter, r *http.Request) {
	rc, err := s.auth.ResolveRequestContext(r)
	if err != nil {
		writeAPIError(w, http.StatusNotFound, "no_active_org", err.Error(), nil)
		return
	}

	// Return organization from context - determined by X-Organization-Slug header
	writeJSON(w, http.StatusOK, CurrentOrganizationResponse{
		Organization: &OrganizationDTO{
			Slug:         rc.OrganizationSlugOrTenant(),
			Name:         rc.OrganizationNameOrTenant(),
			Country:      rc.OrganizationCountryOrTenant(),
			BaseCurrency: rc.OrganizationBaseCurrencyOrTenant(),
		},
	})
}

func (s *Server) requireIdentity(w http.ResponseWriter, r *http.Request) (auth.Identity, bool) {
	identity, err := s.auth.Verify(r)
	if err != nil {
		writeAPIError(w, http.StatusUnauthorized, "unauthenticated", err.Error(), nil)
		return auth.Identity{}, false
	}
	return identity, true
}
