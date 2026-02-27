// package main

// import (
// 	"backend_go/config"
// 	"backend_go/helper"
// 	"backend_go/routes"

// 	"github.com/gin-gonic/gin"
// 	"github.com/joho/godotenv"
// )

// func main() {

// 	godotenv.Load(".env")

// 	r := gin.Default()

// 	r.Use(helper.CorsMiddleware())

// 	db := config.InitDB()

// 	routes.SetupRoutes(r, db)

// 	r.Run(":7070")
// }
// package main

// import (
// 	"backend_go/config"
// 	"backend_go/helper"
// 	"backend_go/routes"

// 	"github.com/gin-gonic/gin"
// 	"github.com/joho/godotenv"
// )

// func main() {

// 	godotenv.Load(".env")

// 	r := gin.Default()

// 	// -----------------------------
// 	// Set trusted proxies di awal
// 	// Ganti sesuai IP proxy kamu, misal Nginx / Load Balancer
// 	// Untuk development lokal bisa pakai nil (percaya semua) tapi ada warning
// 	r.SetTrustedProxies([]string{"127.0.0.1"}) // contoh: localhost proxy
// 	// -----------------------------

// 	r.Use(helper.CorsMiddleware())

// 	db := config.InitDB()

// 	routes.SetupRoutes(r, db)

//		r.Run(":7070")
//	}
// package main

// import (
// 	"backend_go/config"
// 	"backend_go/helper"
// 	"backend_go/routes"
// 	"log"
// 	"net/http"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/joho/godotenv"
// )

// func main() {
// 	// Load env
// 	godotenv.Load(".env")

// 	// Gin mode
// 	gin.SetMode(gin.ReleaseMode)

// 	r := gin.New()

// 	// Middleware wajib
// 	r.Use(
// 		gin.Recovery(),
// 		gin.Logger(),
// 		helper.CorsMiddleware(),
// 	)

// 	// Trusted proxy (sesuaikan environment)
// 	r.SetTrustedProxies([]string{"127.0.0.1"})

// 	// Init DB (connection pool)
// 	db, err := config.InitDB()
// 	if err != nil {
// 		log.Fatal("Failed to connect database:", err)
// 	}

// 	// Setup routes
// 	routes.SetupRoutes(r, db)

// 	r.Static("/qrcodes", "./qrcode")

// 	// HTTP Server dengan timeout
// 	srv := &http.Server{
// 		Addr:         ":7070",
// 		Handler:      r,
// 		ReadTimeout:  5 * time.Second,
// 		WriteTimeout: 10 * time.Second,
// 		IdleTimeout:  120 * time.Second,
// 	}

// 	log.Println("Server running on port 7070")
// 	log.Fatal(srv.ListenAndServe())
// }

package main

import (
	"backend_go/config"
	"backend_go/helper"
	"backend_go/routes"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	godotenv.Load(".env")

	// Set Gin mode to release mode
	gin.SetMode(gin.ReleaseMode)

	// Initialize Gin router
	r := gin.New()

	// Apply mandatory middleware
	r.Use(
		gin.Recovery(),
		gin.Logger(),
		helper.CorsMiddleware(),
	)

	// Set trusted proxies (adjust based on your environment)
	r.SetTrustedProxies([]string{"127.0.0.1"})

	// Initialize database connection
	db, err := config.InitDB()
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	// Setup cron job for auto-expire vouchers
	voucherCron := helper.NewVoucherCronJob(db)
	//scheduler := voucherCron.StartVoucherCron()
	scheduler := voucherCron.StartVoucherCronForTesting()

	// Setup all routes
	routes.SetupRoutes(r, db)

	// Serve static files for QR codes
	r.Static("/qrcodes", "./qrcode")

	// Graceful shutdown handler
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		log.Println("[APP] Shutting down server...")

		// Stop cron scheduler gracefully
		if scheduler != nil {
			scheduler.Stop()
			log.Println("[APP] Cron scheduler stopped")
		}

		os.Exit(0)
	}()

	// Configure HTTP server with timeouts
	srv := &http.Server{
		Addr:         ":7070",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Println("[APP] Server starting on port 7070...")
	log.Println("[APP] Voucher auto-expire cron job is running (every hour)")

	// Start server
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("[APP] Server error:", err)
	}
}
