package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitTicketTemplateRoutes(r *gin.Engine, db *gorm.DB) {
	// Service untuk auth (WAJIB)
	userService := service.NewUserService(db)

	// Service untuk business logic
	ticketTemplateService := service.NewTicketTemplateService(db)
	ticketTemplateController := controller.NewTicketTemplateController(ticketTemplateService)

	ticketTemplateGroup := r.Group("/ticket-templates")
	{
		ticketTemplateGroup.Use(
			middleware.AuthMiddleware(userService), // ✅ FIX
			middleware.RoleMiddleware("admin", "user"),
		)

		ticketTemplateGroup.POST("", ticketTemplateController.CreateTicketTemplate)
		ticketTemplateGroup.GET("", ticketTemplateController.GetAllTicketTemplates)
		ticketTemplateGroup.GET("/:id", ticketTemplateController.GetTicketTemplateByID)
		ticketTemplateGroup.PUT("/:id", ticketTemplateController.UpdateTicketTemplate)
		ticketTemplateGroup.DELETE("/:id", ticketTemplateController.DeleteTicketTemplate)
	}
}
