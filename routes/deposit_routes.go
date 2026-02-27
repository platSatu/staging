package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitDepositRoutes(r *gin.Engine, db *gorm.DB) {
	// Buat userService untuk middleware (konsisten dengan user_routes.go)
	userService := service.NewUserService(db)

	depositService := service.NewDepositService(db)
	depositController := controller.NewDepositController(depositService)

	depositGroup := r.Group("/deposits")
	{
		// Auth middleware menggunakan userService (bukan depositService)
		depositGroup.Use(
			middleware.AuthMiddleware(userService), // Diperbaiki: Gunakan userService seperti pada user
			middleware.RoleMiddleware("admin", "user"),
		)

		depositGroup.POST("", depositController.CreateDeposit)
		depositGroup.GET("", depositController.GetAllDeposits)
		depositGroup.GET("/:id", depositController.GetDepositByID)
		depositGroup.PUT("/:id", depositController.UpdateDeposit)
		depositGroup.DELETE("/:id", depositController.DeleteDeposit)
		// Hapus GET /user karena sudah ada GET "" yang filtered by user
	}
}
