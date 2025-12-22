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
	// Service untuk business logic
	purchaseService := service.NewPurchaseService(db)
	purchaseController := controller.NewPurchaseController(purchaseService)

	purchaseGroup := r.Group("/ticket-purchases")
	{
		purchaseGroup.Use(middleware.RateLimiter(3, 5))
		// Tidak pakai middleware
		purchaseGroup.POST("", purchaseController.ProcessPurchase)
	}

	// Route GET untuk /purchase?token=... (public, tanpa middleware)
	r.GET("/purchase", purchaseController.GetPurchaseByToken)
}
