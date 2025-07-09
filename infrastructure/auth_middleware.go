package infrastructure

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware interface {
	Authenticate() gin.HandlerFunc
}

type authMiddleware struct {
	tokenService TokenService
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(tokenService TokenService) AuthMiddleware {
	return &authMiddleware{tokenService: tokenService}
}

// Authenticate middleware
func (m *authMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		username, err := m.tokenService.ValidateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired session token"})
			return
		}

		c.Set("user", username)
		c.Next()
	}
}
