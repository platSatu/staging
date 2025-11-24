package helper

import (
	"github.com/gin-gonic/gin"
)

// CorsMiddleware adalah middleware untuk menangani CORS
// Hanya menerima request dari origin http://localhost:3000
func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		allowedOrigin := "http://localhost:3000" // SESUAI FRONTEND MU

		if origin == allowedOrigin {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
}
