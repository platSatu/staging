package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitTicketUserRoutes(r *gin.Engine, db *gorm.DB) {
	userService := service.NewUserService(db)
	ticketUserService := service.NewTicketUserService(db)
	ticketUserController := controller.NewTicketUserController(ticketUserService)

	ticketUserGroup := r.Group("/ticket-user")
	{
		ticketUserGroup.Use(
			middleware.AuthMiddleware(userService), // Asumsikan middleware serupa, sesuaikan jika perlu
			middleware.RoleMiddleware("admin", "user"),
		)

		// Endpoint untuk mendapatkan ticket qrcodes milik user yang login
		ticketUserGroup.GET("", ticketUserController.GetMyTicketQrcodes)
	}
}
