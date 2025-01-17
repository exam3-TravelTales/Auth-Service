package middleware

import (
	"auth/api/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Check(c *gin.Context) {
	accessToken := c.GetHeader("Authorization")

	if accessToken == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Authorization is required",
		})
		return
	}

	_, err := auth.ValidateAccessToken(accessToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}
	c.Next()
}
