package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitPaymentPenaltiesRoutes(r *gin.Engine, db *gorm.DB) {
	paymentPenaltiesService := service.NewPaymentPenaltiesService(db)
	paymentPenaltiesController := controller.NewPaymentPenaltiesController(paymentPenaltiesService)

	userService := service.NewUserService(db) // Untuk AuthMiddleware

	paymentPenaltiesGroup := r.Group("/payment-penalties")
	paymentPenaltiesGroup.Use(
		middleware.AuthMiddleware(userService), // Mengambil userID dari token
		middleware.RoleMiddleware("admin", "user"),
	)

	paymentPenaltiesGroup.POST("", paymentPenaltiesController.CreatePaymentPenalties)
	paymentPenaltiesGroup.GET("", paymentPenaltiesController.GetAllPaymentPenalties)
	paymentPenaltiesGroup.GET("/:id", paymentPenaltiesController.GetPaymentPenaltiesByID)
	paymentPenaltiesGroup.PUT("/:id", paymentPenaltiesController.UpdatePaymentPenalties)
	paymentPenaltiesGroup.DELETE("/:id", paymentPenaltiesController.DeletePaymentPenalties)
}
