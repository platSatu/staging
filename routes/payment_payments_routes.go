package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitPaymentPaymentsRoutes(r *gin.Engine, db *gorm.DB) {
	paymentPaymentsService := service.NewPaymentPaymentsService(db)
	paymentPaymentsController := controller.NewPaymentPaymentsController(paymentPaymentsService)

	userService := service.NewUserService(db) // Untuk AuthMiddleware

	paymentPaymentsGroup := r.Group("/payment-payments")
	paymentPaymentsGroup.Use(
		middleware.AuthMiddleware(userService), // Mengambil userID dari token
		middleware.RoleMiddleware("admin", "user"),
	)

	paymentPaymentsGroup.POST("", paymentPaymentsController.CreatePaymentPayments)
	paymentPaymentsGroup.GET("", paymentPaymentsController.GetAllPaymentPayments)
	paymentPaymentsGroup.GET("/:id", paymentPaymentsController.GetPaymentPaymentsByID)
	paymentPaymentsGroup.PUT("/:id", paymentPaymentsController.UpdatePaymentPayments)
	paymentPaymentsGroup.DELETE("/:id", paymentPaymentsController.DeletePaymentPayments)
}
