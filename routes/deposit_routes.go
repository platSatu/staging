package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitDepositRoutes(r *gin.Engine, db *gorm.DB) {
	depositService := service.NewDepositService(db)
	depositController := controller.NewDepositController(depositService)

	depositGroup := r.Group("/deposits")
	{
		// Tambahkan middleware jika perlu, misalnya:
		// depositGroup.Use(middleware.AuthMiddleware(service.NewUserService(db)))

		depositGroup.POST("", depositController.CreateDeposit)
		depositGroup.GET("", depositController.GetAllDeposits)
		depositGroup.GET("/:id", depositController.GetDepositByID)
		depositGroup.GET("/user", depositController.GetDepositsByUser) // Endpoint khusus untuk user
		depositGroup.PUT("/:id", depositController.UpdateDeposit)
		depositGroup.DELETE("/:id", depositController.DeleteDeposit)
	}
}
