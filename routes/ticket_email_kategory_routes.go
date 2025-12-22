package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitTicketEmailKategoryRoutes(r *gin.Engine, db *gorm.DB) {
	// Service untuk auth (WAJIB)
	userService := service.NewUserService(db)

	// Service untuk business logic
	ticketEmailKategoryService := service.NewTicketEmailKategoryService(db)
	ticketEmailKategoryController := controller.NewTicketEmailKategoryController(ticketEmailKategoryService)

	ticketEmailKategoryGroup := r.Group("/ticket-email-kategories")
	{
		ticketEmailKategoryGroup.Use(
			middleware.AuthMiddleware(userService), // ✅ FIX
			middleware.RoleMiddleware("admin", "user"),
		)

		ticketEmailKategoryGroup.POST("", ticketEmailKategoryController.CreateTicketEmailKategory)
		ticketEmailKategoryGroup.GET("", ticketEmailKategoryController.GetAllTicketEmailKategories)
		ticketEmailKategoryGroup.GET("/:id", ticketEmailKategoryController.GetTicketEmailKategoryByID)
		ticketEmailKategoryGroup.PUT("/:id", ticketEmailKategoryController.UpdateTicketEmailKategory)
		ticketEmailKategoryGroup.DELETE("/:id", ticketEmailKategoryController.DeleteTicketEmailKategory)
	}
}
