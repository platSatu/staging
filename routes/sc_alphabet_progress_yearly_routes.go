package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitScAlphabetProgressYearlyRoutes(r *gin.Engine, db *gorm.DB) {
	// Buat userService untuk middleware (konsisten dengan user_routes.go)
	userService := service.NewUserService(db)

	progressService := service.NewScAlphabetProgressYearlyService(db)
	progressController := controller.NewScAlphabetProgressYearlyController(progressService)

	progressGroup := r.Group("/sc-alphabet-progress-yearly")
	{
		// Auth middleware menggunakan userService (bukan progressService)
		progressGroup.Use(
			middleware.AuthMiddleware(userService), // Diperbaiki: Gunakan userService seperti pada user
			middleware.RoleMiddleware("admin", "user"),
		)

		progressGroup.POST("", progressController.CreateScAlphabetProgressYearly)
		progressGroup.GET("", progressController.GetAllScAlphabetProgressYearly)
		progressGroup.GET("/:id", progressController.GetScAlphabetProgressYearlyByID)
		progressGroup.PUT("/:id", progressController.UpdateScAlphabetProgressYearly)
		progressGroup.DELETE("/:id", progressController.DeleteScAlphabetProgressYearly)
	}
}