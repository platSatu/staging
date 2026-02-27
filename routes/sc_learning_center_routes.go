package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitScLearningCenterRoutes(r *gin.Engine, db *gorm.DB) {
	// Buat userService untuk middleware (konsisten dengan user_routes.go)
	userService := service.NewUserService(db)

	learningCenterService := service.NewScLearningCenterService(db)
	learningCenterController := controller.NewScLearningCenterController(learningCenterService)

	learningCenterGroup := r.Group("/sc-learning-center")
	{
		// Auth middleware menggunakan userService (bukan learningCenterService)
		learningCenterGroup.Use(
			middleware.AuthMiddleware(userService), // Diperbaiki: Gunakan userService seperti pada user
			middleware.RoleMiddleware("admin", "user"),
		)

		learningCenterGroup.POST("", learningCenterController.CreateScLearningCenter)
		learningCenterGroup.GET("", learningCenterController.GetAllScLearningCenter)
		learningCenterGroup.GET("/:id", learningCenterController.GetScLearningCenterByID)
		learningCenterGroup.PUT("/:id", learningCenterController.UpdateScLearningCenter)
		learningCenterGroup.DELETE("/:id", learningCenterController.DeleteScLearningCenter)
	}
}
