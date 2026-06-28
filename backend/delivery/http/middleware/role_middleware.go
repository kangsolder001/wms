package middleware

import (
	"net/http"

	"wms/pkg/logger"
)

type RoleMiddleware struct {
	log logger.Logger
}

func NewRoleMiddleware(log logger.Logger) *RoleMiddleware {
	return &RoleMiddleware{log: log}
}

type roleEnforcer struct {
	*RoleMiddleware
	roles []string
}

func (r *roleEnforcer) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		role, ok := req.Context().Value("role").(string)
		if !ok {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		allowed := false
		for _, allowedRole := range r.roles {
			if role == allowedRole {
				allowed = true
				break
			}
		}

		if !allowed {
			r.log.Warn("access denied", "required_roles", r.roles, "user_role", role)
			http.Error(w, `{"error":"insufficient permissions"}`, http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, req)
	})
}

func (m *RoleMiddleware) RequireRole(roles ...string) *roleEnforcer {
	return &roleEnforcer{RoleMiddleware: m, roles: roles}
}
