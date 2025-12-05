package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitPaymentFormRoutes(r *gin.Engine, db *gorm.DB) {
	paymentFormService := service.NewPaymentFormService(db)
	paymentFormController := controller.NewPaymentFormController(paymentFormService)

	userService := service.NewUserService(db) // Untuk AuthMiddleware

	paymentFormGroup := r.Group("/payment-forms")
	paymentFormGroup.Use(
		middleware.AuthMiddleware(userService), // Mengambil userID dari token
		middleware.RoleMiddleware("admin", "user"),
	)

	paymentFormGroup.POST("", paymentFormController.CreatePaymentForm)
	paymentFormGroup.GET("", paymentFormController.GetAllPaymentForms)
	paymentFormGroup.GET("/:id", paymentFormController.GetPaymentFormByID)
	paymentFormGroup.PUT("/:id", paymentFormController.UpdatePaymentForm)
	paymentFormGroup.DELETE("/:id", paymentFormController.DeletePaymentForm)
}
