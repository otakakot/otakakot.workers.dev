package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/syumai/workers/cloudflare"
)

func CORS() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			allowOrigins := parseOrigins(cloudflare.Getenv("CORS_ALLOW_ORIGINS"))

			origin := r.Header.Get("Origin")

			if origin == "" {
				next.ServeHTTP(w, r)
				return
			}

			allowed := matchOrigin(origin, allowOrigins)
			if allowed == "" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				_ = json.NewEncoder(w).Encode(map[string]string{
					"message": "origin not allowed",
				})
				return
			}

			w.Header().Set("Access-Control-Allow-Origin", allowed)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Add("Vary", "Origin")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func parseOrigins(value string) []string {
	if value == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	origins := make([]string, 0, len(parts))
	for _, p := range parts {
		o := strings.TrimSpace(p)
		if o != "" {
			origins = append(origins, o)
		}
	}
	return origins
}

func matchOrigin(origin string, allowed []string) string {
	for _, o := range allowed {
		if o == origin {
			return o
		}
	}
	return ""
}
