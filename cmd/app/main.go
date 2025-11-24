package main

import (
	"backend_go/config"
	"backend_go/helper" // Tambahkan import untuk helper (untuk CorsMiddleware)
	"backend_go/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv" // Tambahkan import ini
)

func main() {
	// Load environment variables dari .env
	godotenv.Load(".env") // Tambahkan ini untuk memuat JWT_SECRET dll.

	r := gin.Default()

	// Tambahkan middleware CORS untuk menerima request dari localhost:3000
	r.Use(helper.CorsMiddleware())

	// Inisialisasi database
	db := config.InitDB() // harus mengembalikan *gorm.DB

	// Setup semua routes, termasuk user
	routes.SetupRoutes(r, db)

	// Jalankan server di port 7070
	r.Run(":7070")
}
