package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitPaymentInstallmentsRoutes(r *gin.Engine, db *gorm.DB) {
	paymentInstallmentsService := service.NewPaymentInstallmentsService(db)
	paymentInstallmentsController := controller.NewPaymentInstallmentsController(paymentInstallmentsService)

	userService := service.NewUserService(db) // Untuk AuthMiddleware

	paymentInstallmentsGroup := r.Group("/payment-installments")
	paymentInstallmentsGroup.Use(
		middleware.AuthMiddleware(userService), // Mengambil userID dari token
		middleware.RoleMiddleware("admin", "user"),
	)

	paymentInstallmentsGroup.POST("", paymentInstallmentsController.CreatePaymentInstallments)
	paymentInstallmentsGroup.GET("", paymentInstallmentsController.GetAllPaymentInstallments)
	paymentInstallmentsGroup.GET("/:id", paymentInstallmentsController.GetPaymentInstallmentsByID)
	paymentInstallmentsGroup.PUT("/:id", paymentInstallmentsController.UpdatePaymentInstallments)
	paymentInstallmentsGroup.DELETE("/:id", paymentInstallmentsController.DeletePaymentInstallments)
}
