package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitTicketResellerSettingRoutes(r *gin.Engine, db *gorm.DB) {
	// Service untuk auth (WAJIB)
	userService := service.NewUserService(db)

	// Service untuk business logic
	ticketResellerSettingService := service.NewTicketResellerSettingService(db)
	ticketResellerSettingController := controller.NewTicketResellerSettingController(ticketResellerSettingService)

	ticketResellerSettingGroup := r.Group("/ticket-reseller-settings")
	{
		ticketResellerSettingGroup.Use(
			middleware.AuthMiddleware(userService), // ✅ FIX
			middleware.RoleMiddleware("admin", "user"),
		)

		ticketResellerSettingGroup.POST("", ticketResellerSettingController.CreateTicketResellerSetting)
		ticketResellerSettingGroup.GET("", ticketResellerSettingController.GetAllTicketResellerSettings)
		ticketResellerSettingGroup.GET("/:id", ticketResellerSettingController.GetTicketResellerSettingByID)
		ticketResellerSettingGroup.PUT("/:id", ticketResellerSettingController.UpdateTicketResellerSetting)
		ticketResellerSettingGroup.DELETE("/:id", ticketResellerSettingController.DeleteTicketResellerSetting)
	}
}
