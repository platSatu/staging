package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitPaymentUserRoutes(r *gin.Engine, db *gorm.DB) {
	paymentUserService := service.NewPaymentUserService(db)
	paymentUserController := controller.NewPaymentUserController(paymentUserService)

	userService := service.NewUserService(db) // Untuk AuthMiddleware

	paymentUserGroup := r.Group("/payment-users")
	paymentUserGroup.Use(
		middleware.AuthMiddleware(userService), // Mengambil userID dari token
		middleware.RoleMiddleware("admin", "user"),
	)

	paymentUserGroup.POST("", paymentUserController.CreatePaymentUser)
	paymentUserGroup.GET("", paymentUserController.GetAllPaymentUsers)
	paymentUserGroup.GET("/:id", paymentUserController.GetPaymentUserByID)
	paymentUserGroup.PUT("/:id", paymentUserController.UpdatePaymentUser)
	paymentUserGroup.DELETE("/:id", paymentUserController.DeletePaymentUser)
}
