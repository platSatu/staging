package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitTicketKategoriRoutes(r *gin.Engine, db *gorm.DB) {
	// Service untuk auth (WAJIB)
	userService := service.NewUserService(db)

	// Service untuk business logic
	ticketKategoriService := service.NewTicketKategoriService(db)
	ticketKategoriController := controller.NewTicketKategoriController(ticketKategoriService)

	r.GET("/public/ticket-kategoris", ticketKategoriController.GetAllTicketKategorisPublic)

	ticketKategoriGroup := r.Group("/ticket-kategoris")
	{
		ticketKategoriGroup.Use(
			middleware.AuthMiddleware(userService), // ✅ FIX
			middleware.RoleMiddleware("admin", "user"),
		)

		ticketKategoriGroup.POST("", ticketKategoriController.CreateTicketKategori)
		ticketKategoriGroup.GET("", ticketKategoriController.GetAllTicketKategoris)
		ticketKategoriGroup.GET("/:id", ticketKategoriController.GetTicketKategoriByID)
		ticketKategoriGroup.PUT("/:id", ticketKategoriController.UpdateTicketKategori)
		ticketKategoriGroup.DELETE("/:id", ticketKategoriController.DeleteTicketKategori)
	}
}
