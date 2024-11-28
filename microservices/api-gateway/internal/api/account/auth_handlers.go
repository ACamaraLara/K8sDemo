package account

import (
	"errors"
	"net/http"

	"github.com/ACamaraLara/K8sBlockChainDemo/shared/jwtManager"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func TokenValidator(jwtManager *jwtManager.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token required"})
			c.Abort()
			return
		}
		// Validate the token
		_, err := jwtManager.ValidateToken(tokenString)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired" + err.Error()})
			} else {
				c.JSON(http.StatusForbidden, gin.H{"error": "Invalid token." + err.Error()})
			}
			c.Abort()
			return
		}
		// Continue to the next handler
		c.Next()
	}
}
