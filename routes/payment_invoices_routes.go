package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitPaymentInvoicesRoutes(r *gin.Engine, db *gorm.DB) {
	paymentInvoicesService := service.NewPaymentInvoicesService(db)
	paymentInvoicesController := controller.NewPaymentInvoicesController(paymentInvoicesService)

	userService := service.NewUserService(db) // Untuk AuthMiddleware

	paymentInvoicesGroup := r.Group("/payment-invoices")
	paymentInvoicesGroup.Use(
		middleware.AuthMiddleware(userService), // Mengambil userID dari token
		middleware.RoleMiddleware("admin", "user"),
	)

	paymentInvoicesGroup.POST("", paymentInvoicesController.CreatePaymentInvoices)
	paymentInvoicesGroup.GET("", paymentInvoicesController.GetAllPaymentInvoices)
	paymentInvoicesGroup.GET("/:id", paymentInvoicesController.GetPaymentInvoicesByID)
	paymentInvoicesGroup.PUT("/:id", paymentInvoicesController.UpdatePaymentInvoices)
	paymentInvoicesGroup.DELETE("/:id", paymentInvoicesController.DeletePaymentInvoices)
}
