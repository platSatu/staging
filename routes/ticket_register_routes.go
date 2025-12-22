package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitTicketRegisterRoutes(r *gin.Engine, db *gorm.DB) {
	// Auth service (WAJIB)
	userService := service.NewUserService(db)

	// Business service
	ticketRegisterService := service.NewTicketRegisterService(db)
	ticketRegisterController := controller.NewTicketRegisterController(ticketRegisterService)

	ticketRegisterGroup := r.Group("/ticket-registers")
	{
		ticketRegisterGroup.Use(
			middleware.AuthMiddleware(userService), // ✅ FIX
			middleware.RoleMiddleware("admin", "user"),
		)

		ticketRegisterGroup.POST("", ticketRegisterController.CreateTicketRegister)
		ticketRegisterGroup.GET("", ticketRegisterController.GetAllTicketRegisters)
		ticketRegisterGroup.GET("/:id", ticketRegisterController.GetTicketRegisterByID)
		ticketRegisterGroup.PUT("/:id", ticketRegisterController.UpdateTicketRegister)
		ticketRegisterGroup.DELETE("/:id", ticketRegisterController.DeleteTicketRegister)
	}

	// PUBLIC (tanpa auth)
	// r.POST("/public-ticket-registers", ticketRegisterController.CreatePublicTicketRegister)
	// r.GET("/purchase/validate", ticketRegisterController.ValidatePurchaseToken)
	r.POST("/public-ticket-registers",
		middleware.RateLimiter(5, 10), // batasi spam
		ticketRegisterController.CreatePublicTicketRegister,
	)

	r.GET("/purchase/validate",
		middleware.RateLimiter(10, 20), // lebih longgar untuk validate
		ticketRegisterController.ValidatePurchaseToken,
	)

}
