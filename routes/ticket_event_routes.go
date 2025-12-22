package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitTicketEventRoutes(r *gin.Engine, db *gorm.DB) {
	// Service untuk auth (WAJIB)
	userService := service.NewUserService(db)

	// Service untuk business logic
	ticketEventService := service.NewTicketEventService(db)
	ticketEventController := controller.NewTicketEventController(ticketEventService)

	r.GET("/public/ticket-events", ticketEventController.GetAllTicketEventsPublic)

	ticketEventGroup := r.Group("/ticket-events")
	{
		ticketEventGroup.Use(
			middleware.AuthMiddleware(userService), // ✅ FIX
			middleware.RoleMiddleware("admin", "user"),
		)

		ticketEventGroup.POST("", ticketEventController.CreateTicketEvent)
		ticketEventGroup.GET("", ticketEventController.GetAllTicketEvents)
		ticketEventGroup.GET("/:id", ticketEventController.GetTicketEventByID)
		ticketEventGroup.PUT("/:id", ticketEventController.UpdateTicketEvent)
		ticketEventGroup.DELETE("/:id", ticketEventController.DeleteTicketEvent)
	}
}
