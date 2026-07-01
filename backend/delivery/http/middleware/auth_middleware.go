package middleware

import (
	"database/sql"
	"net/http"
	"strings"

	"wms/infrastructure/auth"
	"wms/pkg/logger"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	jwtService auth.JWTService
	db         *sql.DB
	log        logger.Logger
}

func NewAuthMiddleware(jwtService auth.JWTService, db *sql.DB, log logger.Logger) *AuthMiddleware {
	return &AuthMiddleware{jwtService: jwtService, db: db, log: log}
}

func (m *AuthMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			return
		}

		claims, err := m.jwtService.ValidateToken(parts[1])
		if err != nil {
			m.log.Error("invalid token", "error", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			return
		}

		var isActive bool
		err = m.db.QueryRow("SELECT is_active FROM users WHERE id = $1", userID).Scan(&isActive)
		if err != nil || !isActive {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "account is deactivated"})
			return
		}

		role, _ := claims["role"].(string)

		c.Set("user_id", userID)
		c.Set("role", role)

		c.Next()
	}
}
