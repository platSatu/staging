// // purchase_routes.go
// package routes

// import (
// 	"backend_go/internal/controller"
// 	"backend_go/internal/service"
// 	"backend_go/middleware"

// 	"github.com/gin-gonic/gin"
// 	"gorm.io/gorm"
// )

// func InitPurchaseRoutes(r *gin.Engine, db *gorm.DB) {
// 	// Service untuk business logic
// 	purchaseService := service.NewPurchaseService(db)
// 	purchaseController := controller.NewPurchaseController(purchaseService)

// 	purchaseGroup := r.Group("/ticket-purchases")
// 	{
// 		purchaseGroup.Use(middleware.RateLimiter(3, 5))
// 		// Tidak pakai middleware
// 		purchaseGroup.POST("", purchaseController.ProcessPurchase)
// 	}

//		// Route GET untuk /purchase?token=... (public, tanpa middleware)
//		// r.GET("/purchase", purchaseController.GetPurchaseByToken)
//	}
//
// purchase_routes.go
package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitPurchaseRoutes(r *gin.Engine, db *gorm.DB) {
	// Service untuk auth (WAJIB, seperti di ticket-kategori)
	userService := service.NewUserService(db)

	// Service untuk business logic
	purchaseService := service.NewPurchaseService(db)
	purchaseController := controller.NewPurchaseController(purchaseService)

	purchaseGroup := r.Group("/ticket-purchases")
	{
		// Tambahkan middleware auth dan role seperti di ticket-kategori
		purchaseGroup.Use(
			middleware.AuthMiddleware(userService),     // ✅ FIX: Tambahkan ini untuk set userID di context
			middleware.RoleMiddleware("admin", "user"), // ✅ FIX: Tambahkan ini untuk cek role
		)
		purchaseGroup.Use(middleware.RateLimiter(3, 5)) // Rate limiter tetap ada
		purchaseGroup.POST("", purchaseController.ProcessPurchase)
	}

	// Route GET untuk /purchase?token=... (public, tanpa middleware) - Tetap seperti sebelumnya jika diperlukan
	// r.GET("/purchase", purchaseController.GetPurchaseByToken)
}
