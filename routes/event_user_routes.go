package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitEventUserRoutes(r *gin.Engine, db *gorm.DB) {
	eventUserService := service.NewEventUserService(db)
	eventUserController := controller.NewEventUserController(eventUserService)

	eventUserGroup := r.Group("/event_users")
	{
		// Tanpa middleware, akses langsung
		eventUserGroup.POST("/register", eventUserController.RegisterEventUser)
	}
}
