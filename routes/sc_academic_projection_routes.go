package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitScAcademicProjectionRoutes(r *gin.Engine, db *gorm.DB) {
	// Buat userService untuk middleware (konsisten dengan user_routes.go)
	userService := service.NewUserService(db)

	projectionService := service.NewScAcademicProjectionService(db)
	projectionController := controller.NewScAcademicProjectionController(projectionService)

	projectionGroup := r.Group("/sc-academic-projection")
	{
		// Auth middleware menggunakan userService (bukan projectionService)
		projectionGroup.Use(
			middleware.AuthMiddleware(userService), // Diperbaiki: Gunakan userService seperti pada user
			middleware.RoleMiddleware("admin", "user"),
		)
		projectionGroup.GET("/user", projectionController.GetScAcademicProjectionByUser)
		projectionGroup.POST("", projectionController.CreateScAcademicProjection)
		projectionGroup.GET("", projectionController.GetAllScAcademicProjection)
		projectionGroup.GET("/:id", projectionController.GetScAcademicProjectionByID)
		projectionGroup.PUT("/:id", projectionController.UpdateScAcademicProjection)
		projectionGroup.DELETE("/:id", projectionController.DeleteScAcademicProjection)
	}
}
