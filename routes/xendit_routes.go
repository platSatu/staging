package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitXenditRoutes(r *gin.Engine, db *gorm.DB) {
	// Setup services
	transaksiService := service.NewTransaksiService(db)
	xenditService := service.NewXenditService(transaksiService)
	xenditController := controller.NewXenditController(xenditService)

	// User service untuk middleware
	userService := service.NewUserService(db)

	// Public route - callback dari Xendit (tidak butuh auth)
	r.POST("/api/xendit/callback", xenditController.Callback)

	// Protected routes - butuh auth
	xenditGroup := r.Group("/api/xendit")
	xenditGroup.Use(
		middleware.AuthMiddleware(userService),
		middleware.RoleMiddleware("admin", "user"),
	)
	{
		xenditGroup.POST("/create-payment", xenditController.CreatePayment)
	}
}
