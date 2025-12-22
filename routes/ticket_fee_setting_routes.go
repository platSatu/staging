package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitTicketFeeSettingRoutes(r *gin.Engine, db *gorm.DB) {
	// Service untuk auth (WAJIB)
	userService := service.NewUserService(db)

	// Service untuk business logic
	ticketFeeSettingService := service.NewTicketFeeSettingService(db)
	ticketFeeSettingController := controller.NewTicketFeeSettingController(ticketFeeSettingService)

	r.GET("/public/ticket-fees", ticketFeeSettingController.GetAllTicketFeeSettingsPublic)

	ticketFeeSettingGroup := r.Group("/ticket-fee-settings")
	{
		ticketFeeSettingGroup.Use(
			middleware.AuthMiddleware(userService), // ✅ FIX
			middleware.RoleMiddleware("admin", "user"),
		)

		ticketFeeSettingGroup.POST("", ticketFeeSettingController.CreateTicketFeeSetting)
		ticketFeeSettingGroup.GET("", ticketFeeSettingController.GetAllTicketFeeSettings)
		ticketFeeSettingGroup.GET("/:id", ticketFeeSettingController.GetTicketFeeSettingByID)
		ticketFeeSettingGroup.PUT("/:id", ticketFeeSettingController.UpdateTicketFeeSetting)
		ticketFeeSettingGroup.DELETE("/:id", ticketFeeSettingController.DeleteTicketFeeSetting)
	}
}
