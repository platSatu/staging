package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitPaymentCategoryRoutes(r *gin.Engine, db *gorm.DB) {
	paymentCategoryService := service.NewPaymentCategoryService(db)
	paymentCategoryController := controller.NewPaymentCategoryController(paymentCategoryService)

	userService := service.NewUserService(db) // Untuk AuthMiddleware

	paymentCategoryGroup := r.Group("/payment-categories")
	paymentCategoryGroup.Use(
		middleware.AuthMiddleware(userService), // Mengambil userID dari token
		middleware.RoleMiddleware("admin", "user"),
	)

	paymentCategoryGroup.POST("", paymentCategoryController.CreatePaymentCategory)
	paymentCategoryGroup.GET("", paymentCategoryController.GetAllPaymentCategories)
	paymentCategoryGroup.GET("/:id", paymentCategoryController.GetPaymentCategoryByID)
	paymentCategoryGroup.PUT("/:id", paymentCategoryController.UpdatePaymentCategory)
	paymentCategoryGroup.DELETE("/:id", paymentCategoryController.DeletePaymentCategory)
}
