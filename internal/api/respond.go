package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

const (
	headerOrganization = "X-Organization-Slug"
	headerTotalCount   = "X-Total-Count"
	headerLimit        = "X-Limit"
	headerOffset       = "X-Offset"
)

// Pure JSON responses - no cookies, no CSRF
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func writeAPIError(w http.ResponseWriter, status int, code, message string, details any) {
	writeJSON(w, status, APIError{
		Code:    code,
		Message: message,
		Details: details,
	})
}

func parseJSON[T any](r *http.Request) (T, error) {
	var out T
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&out)
	return out, err
}

func parseListQuery(r *http.Request) ListMeta {
	return ListMeta{
		Limit:  clampInt(queryInt(r, "limit", 50), 1, 200),
		Offset: clampInt(queryInt(r, "offset", 0), 0, 1_000_000),
	}
}

func queryString(r *http.Request, key string) string {
	return strings.TrimSpace(r.URL.Query().Get(key))
}

func queryInt(r *http.Request, key string, fallback int) int {
	v := r.URL.Query().Get(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}

func clampInt(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
