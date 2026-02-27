package helper

import (
	"log"

	"github.com/gin-gonic/gin"
)

// CorsMiddleware adalah middleware untuk menangani CORS
// Hanya menerima request dari origin http://localhost:3000
// func CorsMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		origin := c.GetHeader("Origin")

// 		allowedOrigin := "http://localhost:3000" // SESUAI FRONTEND MU

// 		if origin == allowedOrigin {
// 			c.Header("Access-Control-Allow-Origin", origin)
// 		}

// 		c.Header("Access-Control-Allow-Credentials", "true")
// 		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
// 		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

// 		if c.Request.Method == "OPTIONS" {
// 			c.AbortWithStatus(200)
// 			return
// 		}

// 		c.Next()
// 	}
// }

// CorsMiddleware adalah middleware untuk menangani CORS
// Mengizinkan origins tertentu dan menangani preflight OPTIONS
func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// Log origin untuk debugging (hapus di production)
		log.Printf("Received Origin: %s", origin)

		// Daftar origin yang diizinkan
		allowedOrigins := map[string]bool{
			"https://event.gbigatsu.id":   true,
			"https://backend.gbigatsu.id": true, // Frontend production
			"http://localhost:3000":       true, // Untuk development
			// Tambahkan origins lain jika perlu, misalnya subdomain
		}

		// Selalu atur header umum
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Jika origin diizinkan, atur header spesifik
		if allowedOrigins[origin] {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true") // Set ke "false" jika tidak perlu credentials
		}

		// Tangani preflight OPTIONS request
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
}
