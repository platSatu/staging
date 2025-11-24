package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitEventRoutes(r *gin.Engine, db *gorm.DB) {
	eventService := service.NewEventService(db)
	eventController := controller.NewEventController(eventService)

	eventGroup := r.Group("/events")
	{
		eventGroup.Use(
			middleware.AuthMiddleware(service.NewUserService(db)), // Menggunakan UserService untuk auth
			middleware.RoleMiddleware("admin", "user"),
		)

		eventGroup.POST("", eventController.CreateEvent)
		eventGroup.GET("", eventController.GetAllEvents)
		eventGroup.GET("/:id", eventController.GetEventByID)
		eventGroup.GET("/user", eventController.GetEventsByUser) // Endpoint khusus untuk user login
		eventGroup.PUT("/:id", eventController.UpdateEvent)
		eventGroup.DELETE("/:id", eventController.DeleteEvent)
	}
}
