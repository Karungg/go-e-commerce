package middleware

import (
	"net/http"
	"strings"

	"go-e-commerce/internal/delivery/http/response"
	"go-e-commerce/internal/port"

	"github.com/gin-gonic/gin"
)

// RequireAuth is a middleware that validates the JWT token
func RequireAuth(jwtAuth port.TokenValidator) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "Authorization failed", "Authorization header is missing")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			response.Error(c, http.StatusUnauthorized, "Authorization failed", "Invalid authorization header format")
			c.Abort()
			return
		}

		claims, err := jwtAuth.ValidateToken(parts[1])
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "Authorization failed", "Invalid or expired token")
			c.Abort()
			return
		}

		// Set user identity to context
		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// RequireRole is a middleware that enforces role-based access control
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			response.Error(c, http.StatusUnauthorized, "Authorization failed", "Role information not found in context")
			c.Abort()
			return
		}

		hasRole := false
		for _, role := range roles {
			if role == userRole {
				hasRole = true
				break
			}
		}

		if !hasRole {
			response.Error(c, http.StatusForbidden, "Forbidden", "Insufficient permissions to access this resource")
			c.Abort()
			return
		}

		c.Next()
	}
}
