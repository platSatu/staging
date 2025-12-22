package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitTicketMethodPembayaranRoutes(r *gin.Engine, db *gorm.DB) {
	// Auth butuh UserService
	userService := service.NewUserService(db)

	// Business logic service
	ticketMethodPembayaranService := service.NewTicketMethodPembayaranService(db)
	ticketMethodPembayaranController :=
		controller.NewTicketMethodPembayaranController(ticketMethodPembayaranService)

	ticketMethodPembayaranGroup := r.Group("/ticket-method-pembayarans")
	{
		ticketMethodPembayaranGroup.Use(
			middleware.AuthMiddleware(userService), // ✅ FIX UTAMA
			middleware.RoleMiddleware("admin", "user"),
		)

		ticketMethodPembayaranGroup.POST("", ticketMethodPembayaranController.CreateTicketMethodPembayaran)
		ticketMethodPembayaranGroup.GET("", ticketMethodPembayaranController.GetAllTicketMethodPembayarans)
		ticketMethodPembayaranGroup.GET("/:id", ticketMethodPembayaranController.GetTicketMethodPembayaranByID)
		ticketMethodPembayaranGroup.PUT("/:id", ticketMethodPembayaranController.UpdateTicketMethodPembayaran)
		ticketMethodPembayaranGroup.DELETE("/:id", ticketMethodPembayaranController.DeleteTicketMethodPembayaran)
	}
}
