package app

import "os"

type Config struct {
	Addr     string
	JWKSURL  string
	JWTIssuer string
}

func LoadConfig() Config {
	return Config{
		Addr:      getEnv("ADDR", ":8080"),
		JWKSURL:   getEnv("JWKS_URL", ""),
		JWTIssuer: getEnv("JWT_ISSUER", ""),
	}
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
