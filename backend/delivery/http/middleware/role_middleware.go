package middleware

import (
	"net/http"

	"wms/pkg/logger"

	"github.com/gin-gonic/gin"
)

type RoleMiddleware struct {
	log logger.Logger
}

func NewRoleMiddleware(log logger.Logger) *RoleMiddleware {
	return &RoleMiddleware{log: log}
}

func (m *RoleMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		roleStr, ok := role.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		allowed := false
		for _, allowedRole := range roles {
			if roleStr == allowedRole {
				allowed = true
				break
			}
		}

		if !allowed {
			m.log.Warn("access denied", "required_roles", roles, "user_role", roleStr)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			return
		}

		c.Next()
	}
}
