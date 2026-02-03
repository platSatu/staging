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
	profileGroup.Use(middleware.AuthMiddleware(service.NewUserService(db)))
	{
		// Admin only
		adminGroup := profileGroup.Group("")
		adminGroup.Use(middleware.RoleMiddleware("admin"))
		{
			adminGroup.GET("", profileController.GetAllProfiles) // Hanya admin
		}

		// User dan admin
		userGroup := profileGroup.Group("")
		userGroup.Use(middleware.RoleMiddleware("admin", "user"))
		{
			userGroup.POST("", profileController.CreateProfile)
			userGroup.GET("/:id", profileController.GetProfileByID)
			userGroup.GET("/user", profileController.GetProfilesByUser)
			userGroup.PUT("/:id", profileController.UpdateProfile)
			userGroup.DELETE("/:id", profileController.DeleteProfile)

			// Tambahkan di dalam profileGroup
			profileGroup.GET("/check", profileController.CheckUserHasProfile) // Bisa diakses user biasa
		}
	}
}
