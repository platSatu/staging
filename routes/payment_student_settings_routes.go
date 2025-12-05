package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitPaymentInvoiceRoutes(r *gin.Engine, db *gorm.DB) {
	paymentInvoiceService := service.NewPaymentInvoiceService(db)
	paymentInvoiceController := controller.NewPaymentInvoiceController(paymentInvoiceService)

	userService := service.NewUserService(db) // Untuk AuthMiddleware

	paymentInvoiceGroup := r.Group("/payment-invoices")
	paymentInvoiceGroup.Use(
		middleware.AuthMiddleware(userService), // Mengambil userID dari token
		middleware.RoleMiddleware("admin", "user"),
	)

	paymentInvoiceGroup.POST("", paymentInvoiceController.CreatePaymentInvoice)
	paymentInvoiceGroup.GET("", paymentInvoiceController.GetAllPaymentInvoices)
	paymentInvoiceGroup.GET("/:id", paymentInvoiceController.GetPaymentInvoiceByID)
	paymentInvoiceGroup.PUT("/:id", paymentInvoiceController.UpdatePaymentInvoice)
	paymentInvoiceGroup.DELETE("/:id", paymentInvoiceController.DeletePaymentInvoice)
}
