package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitScAlphabetProgressRoutes(r *gin.Engine, db *gorm.DB) {
	// Buat userService untuk middleware (konsisten dengan user_routes.go)
	userService := service.NewUserService(db)

	progressService := service.NewScAlphabetProgressService(db)
	progressController := controller.NewScAlphabetProgressController(progressService)

	progressGroup := r.Group("/sc-alphabet-progress")
	{
		// Auth middleware menggunakan userService (bukan progressService)
		progressGroup.Use(
			middleware.AuthMiddleware(userService), // Diperbaiki: Gunakan userService seperti pada user
			middleware.RoleMiddleware("admin", "user"),
		)

		progressGroup.POST("", progressController.CreateScAlphabetProgress)
		progressGroup.GET("", progressController.GetAllScAlphabetProgress)
		progressGroup.GET("/:id", progressController.GetScAlphabetProgressByID)
		progressGroup.PUT("/:id", progressController.UpdateScAlphabetProgress)
		progressGroup.DELETE("/:id", progressController.DeleteScAlphabetProgress)
	}
}