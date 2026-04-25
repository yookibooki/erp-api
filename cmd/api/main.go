package main

import (
	"log"
	"net/http"
	"os"

	"github.com/yookibooki/erp-api/internal/api"
	"github.com/yookibooki/erp-api/internal/app"
	"github.com/yookibooki/erp-api/internal/auth"
)

func main() {
	cfg := app.LoadConfig()
	
	verifier, err := auth.NewVerifier(cfg.JWKSURL, cfg.JWTIssuer)
	if err != nil {
		log.Fatalf("auth: %v", err)
	}

	// TODO: wire real services (db, organizations, etc.)
	apiSrv := api.NewServer(api.Dependencies{
		Auth: verifier,
		OrganizationsService: nil, // inject real service
	})

	addr := cfg.Addr
	if addr == "" {
		addr = ":8080"
	}
	log.Printf("erp-api listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, apiSrv.Router()))
}
