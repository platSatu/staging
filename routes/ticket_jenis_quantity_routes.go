package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitTicketJenisQuantityRoutes(r *gin.Engine, db *gorm.DB) {
	// Service untuk auth (WAJIB)
	userService := service.NewUserService(db)

	// Service untuk business logic
	ticketJenisQuantityService := service.NewTicketJenisQuantityService(db)
	ticketJenisQuantityController := controller.NewTicketJenisQuantityController(ticketJenisQuantityService)

	r.GET("/public/ticket-jenis-quantities", ticketJenisQuantityController.GetAllTicketJenisQuantitiesPublic)

	ticketJenisQuantityGroup := r.Group("/ticket-jenis-quantities")
	{
		ticketJenisQuantityGroup.Use(
			middleware.AuthMiddleware(userService), // ✅ FIX
			middleware.RoleMiddleware("admin", "user"),
		)

		ticketJenisQuantityGroup.POST("", ticketJenisQuantityController.CreateTicketJenisQuantity)
		ticketJenisQuantityGroup.GET("", ticketJenisQuantityController.GetAllTicketJenisQuantities)
		ticketJenisQuantityGroup.GET("/:id", ticketJenisQuantityController.GetTicketJenisQuantityByID)
		ticketJenisQuantityGroup.PUT("/:id", ticketJenisQuantityController.UpdateTicketJenisQuantity)
		ticketJenisQuantityGroup.DELETE("/:id", ticketJenisQuantityController.DeleteTicketJenisQuantity)
	}
}
