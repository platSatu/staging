package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitTicketBlastRoutes(r *gin.Engine, db *gorm.DB) {
	// Service untuk auth (WAJIB)
	userService := service.NewUserService(db)

	// Service untuk business logic
	ticketBlastService := service.NewTicketBlastService(db)
	ticketBlastController := controller.NewTicketBlastController(ticketBlastService)

	ticketBlastGroup := r.Group("/ticket-blasts")
	{
		ticketBlastGroup.Use(
			middleware.AuthMiddleware(userService), // ✅ FIX
			middleware.RoleMiddleware("admin", "user"),
		)

		ticketBlastGroup.POST("", ticketBlastController.CreateTicketBlast)
		ticketBlastGroup.GET("", ticketBlastController.GetAllTicketBlasts)
		ticketBlastGroup.GET("/:id", ticketBlastController.GetTicketBlastByID)
		ticketBlastGroup.PUT("/:id", ticketBlastController.UpdateTicketBlast)
		ticketBlastGroup.DELETE("/:id", ticketBlastController.DeleteTicketBlast)
	}
}
