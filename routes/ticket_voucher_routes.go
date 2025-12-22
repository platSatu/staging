package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitTicketVoucherRoutes(r *gin.Engine, db *gorm.DB) {
	// Service untuk auth (WAJIB)
	userService := service.NewUserService(db)

	// Service untuk business logic
	ticketVoucherService := service.NewTicketVoucherService(db)
	ticketVoucherController := controller.NewTicketVoucherController(ticketVoucherService)

	r.GET("/ticket-vouchers/check/:kodeVoucher", ticketVoucherController.CheckVoucherAvailability)

	ticketVoucherGroup := r.Group("/ticket-vouchers")
	{
		ticketVoucherGroup.Use(
			middleware.AuthMiddleware(userService), // ✅ FIX
			middleware.RoleMiddleware("admin", "user"),
		)

		ticketVoucherGroup.POST("", ticketVoucherController.CreateTicketVoucher)
		ticketVoucherGroup.GET("", ticketVoucherController.GetAllTicketVouchers)
		ticketVoucherGroup.GET("/:id", ticketVoucherController.GetTicketVoucherByID)
		ticketVoucherGroup.PUT("/:id", ticketVoucherController.UpdateTicketVoucher)
		ticketVoucherGroup.DELETE("/:id", ticketVoucherController.DeleteTicketVoucher)
	}
}
