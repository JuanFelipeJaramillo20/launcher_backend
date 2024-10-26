package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"venecraft-back/cmd/utils"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func GetLoggedInUser(c *gin.Context) (uint64, []string, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, nil, false
	}
	roles, _ := c.Get("role")
	return userID.(uint64), roles.([]string), true
}
