package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitTicketQrcodeRoutes(r *gin.Engine, db *gorm.DB) {
	userService := service.NewUserService(db) 
	ticketQrcodeService := service.NewTicketQrcodeService(db)
	ticketQrcodeController := controller.NewTicketQrcodeController(ticketQrcodeService)

	ticketQrcodeGroup := r.Group("/ticket-qrcodes")
	{
		ticketQrcodeGroup.Use(
			middleware.AuthMiddleware(userService), // Asumsikan middleware serupa, sesuaikan jika perlu
			middleware.RoleMiddleware("admin", "user"),
		)

		ticketQrcodeGroup.POST("", ticketQrcodeController.CreateTicketQrcode)
		ticketQrcodeGroup.GET("", ticketQrcodeController.GetAllTicketQrcodes)
		ticketQrcodeGroup.GET("/:id", ticketQrcodeController.GetTicketQrcodeByID)
		ticketQrcodeGroup.PUT("/:id", ticketQrcodeController.UpdateTicketQrcode)
		ticketQrcodeGroup.DELETE("/:id", ticketQrcodeController.DeleteTicketQrcode)
	}
}