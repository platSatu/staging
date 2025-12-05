package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitProfileRoutes(r *gin.Engine, db *gorm.DB) {
	profileService := service.NewProfileService(db)
	profileController := controller.NewProfileController(profileService)

	profileGroup := r.Group("/profiles")
	{
		profileGroup.Use(
			middleware.AuthMiddleware(service.NewUserService(db)), // Menggunakan UserService untuk auth
			middleware.RoleMiddleware("admin", "user"),
		)

		profileGroup.POST("", profileController.CreateProfile)
		profileGroup.GET("", profileController.GetAllProfiles)
		profileGroup.GET("/:id", profileController.GetProfileByID)
		profileGroup.GET("/user", profileController.GetProfilesByUser) // Endpoint khusus untuk user login
		profileGroup.PUT("/:id", profileController.UpdateProfile)
		profileGroup.DELETE("/:id", profileController.DeleteProfile)
	}
}
