package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sh1ro/todo-api/pkg/config"
)

// CORS returns a middleware that handles CORS
func CORS(cfg config.CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set CORS headers
		c.Writer.Header().Set("Access-Control-Allow-Origin", joinStrings(cfg.AllowedOrigins))
		c.Writer.Header().Set("Access-Control-Allow-Methods", joinStrings(cfg.AllowedMethods))
		c.Writer.Header().Set("Access-Control-Allow-Headers", joinStrings(cfg.AllowedHeaders))
		c.Writer.Header().Set("Access-Control-Max-Age", fmt.Sprintf("%d", cfg.MaxAge))
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// joinStrings joins a slice of strings with a comma
func joinStrings(strings []string) string {
	if len(strings) == 0 {
		return ""
	}

	if len(strings) == 1 && strings[0] == "*" {
		return "*"
	}

	result := strings[0]
	for i := 1; i < len(strings); i++ {
		result += ", " + strings[i]
	}

	return result
}
