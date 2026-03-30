package middleware

import (
	"crypto/subtle"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const HeaderAPIKey = "X-API-Key"

func RequireAPIKey(expectedKey string) gin.HandlerFunc {
	expected := strings.TrimSpace(expectedKey)
	return func(c *gin.Context) {
		if expected == "" {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
				"error": "PUBLIC_API_KEY não configurada",
			})
			return
		}
		got := strings.TrimSpace(c.GetHeader(HeaderAPIKey))
		if len(got) != len(expected) || subtle.ConstantTimeCompare([]byte(got), []byte(expected)) != 1 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "API key inválida ou ausente",
			})
			return
		}
		c.Next()
	}
}
