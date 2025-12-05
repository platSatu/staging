package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitPaymentPenaltySettingsRoutes(r *gin.Engine, db *gorm.DB) {
	paymentPenaltySettingsService := service.NewPaymentPenaltySettingsService(db)
	paymentPenaltySettingsController := controller.NewPaymentPenaltySettingsController(paymentPenaltySettingsService)

	userService := service.NewUserService(db) // Untuk AuthMiddleware

	paymentPenaltySettingsGroup := r.Group("/payment-penalty-settings")
	paymentPenaltySettingsGroup.Use(
		middleware.AuthMiddleware(userService), // Mengambil userID dari token
		middleware.RoleMiddleware("admin", "user"),
	)

	paymentPenaltySettingsGroup.POST("", paymentPenaltySettingsController.CreatePaymentPenaltySettings)
	paymentPenaltySettingsGroup.GET("", paymentPenaltySettingsController.GetAllPaymentPenaltySettings)
	paymentPenaltySettingsGroup.GET("/:id", paymentPenaltySettingsController.GetPaymentPenaltySettingsByID)
	paymentPenaltySettingsGroup.PUT("/:id", paymentPenaltySettingsController.UpdatePaymentPenaltySettings)
	paymentPenaltySettingsGroup.DELETE("/:id", paymentPenaltySettingsController.DeletePaymentPenaltySettings)
}
