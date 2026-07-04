package httpserver

import (
	"context"
	"net/http"
	"strings"
	"vincehpicton/click/internal/tokens"
)

func (s *Server) middlewareExample(prevHandler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		prevHandler(w, r)
	}
}

// CORSMiddleware wraps the entire mux router so every response (including 404/405
// from gorilla/mux) includes CORS headers. Router.Use only runs after a route match.
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type AuthContext struct {
	UserID string
}
type contextKey string

const AuthContextKey contextKey = "auth"

func AuthFromContext(ctx context.Context) (AuthContext, bool) {
	v, ok := ctx.Value(AuthContextKey).(AuthContext)
	return v, ok
}

func AuthMiddleware(tokenMgr *tokens.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 {
				http.Error(w, "invalid authorization scheme", http.StatusUnauthorized)
				return
			}
			bearerStr := parts[0]
			tokenString := parts[1]

			if bearerStr != "Bearer" && bearerStr != "bearer" {
				http.Error(w, "invalid auth header prefix; 'Bearer ' or 'bearer' is required", http.StatusUnauthorized)
				return
			}

			claims, err := tokenMgr.ParseAccessToken(tokenString)
			if err != nil {
				http.Error(w, "failed to parse token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), AuthContextKey, AuthContext{
				UserID: claims.Subject,
			})

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
