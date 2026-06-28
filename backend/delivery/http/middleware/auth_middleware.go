package middleware

import (
	"context"
	"net/http"
	"strings"

	"wms/infrastructure/auth"
	"wms/pkg/logger"
)

type AuthMiddleware struct {
	jwtService auth.JWTService
	log        logger.Logger
}

func NewAuthMiddleware(jwtService auth.JWTService, log logger.Logger) *AuthMiddleware {
	return &AuthMiddleware{jwtService: jwtService, log: log}
}

func (m *AuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error":"missing authorization header"}`, http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, `{"error":"invalid authorization header format"}`, http.StatusUnauthorized)
			return
		}

		claims, err := m.jwtService.ValidateToken(parts[1])
		if err != nil {
			m.log.Error("invalid token", "error", err)
			http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			http.Error(w, `{"error":"invalid token claims"}`, http.StatusUnauthorized)
			return
		}

		role, _ := claims["role"].(string)

		ctx := context.WithValue(r.Context(), "user_id", userID)
		ctx = context.WithValue(ctx, "role", role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
