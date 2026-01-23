package middleware

import (
	"gin-learning/pkg/auth"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("auth_token")
		if err != nil {
			slog.Warn("Unauthorized access: cookie not found", "Client IP =", c.ClientIP())
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Please login"})
			return
		}

		claims, err := auth.VerifyToken(token)
		if err != nil {
			slog.Warn("Unauthorized access: invalid token", "error", err, "Client IP =", c.ClientIP())
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
			return
		}

		slog.Info("User token check and it was ok", "Client IP =", c.ClientIP())
		c.Set("userID", claims.UserID)
		c.Set("userLevel", claims.Level)
		c.Next()
	}
}

func AuthorizeRole(allowedLevels ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userLevel, exists := c.Get("userLevel")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Level not found"})
			slog.Warn("Level not found", "Client IP =", c.ClientIP())
			return
		}

		levelStr := userLevel.(string)
		allowed := false
		for _, l := range allowedLevels {
			if l == levelStr {
				allowed = true
				break
			}
		}

		if !allowed {
			slog.Warn("Access denied", "required", allowedLevels, "actual", levelStr)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
		c.Next()
	}
}
