package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitTicketHistoryRoutes(r *gin.Engine, db *gorm.DB) {
	// Service untuk auth (WAJIB)
	userService := service.NewUserService(db)

	// Service untuk business logic
	ticketHistoryService := service.NewTicketHistoryService(db)
	ticketHistoryController := controller.NewTicketHistoryController(ticketHistoryService)

	ticketHistoryGroup := r.Group("/ticket-histories")
	{
		ticketHistoryGroup.Use(
			middleware.AuthMiddleware(userService), // ✅ FIX
			middleware.RoleMiddleware("admin", "user"),
		)

		ticketHistoryGroup.POST("", ticketHistoryController.CreateTicketHistory)
		ticketHistoryGroup.GET("", ticketHistoryController.GetAllTicketHistories)
		ticketHistoryGroup.GET("/:id", ticketHistoryController.GetTicketHistoryByID)
		ticketHistoryGroup.PUT("/:id", ticketHistoryController.UpdateTicketHistory)
		ticketHistoryGroup.DELETE("/:id", ticketHistoryController.DeleteTicketHistory)
	}
}
