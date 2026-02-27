package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitScAcademicProjectionDetailRoutes(r *gin.Engine, db *gorm.DB) {
	// Buat userService untuk middleware (konsisten dengan user_routes.go)
	userService := service.NewUserService(db)

	detailService := service.NewScAcademicProjectionDetailService(db)
	detailController := controller.NewScAcademicProjectionDetailController(detailService)

	detailGroup := r.Group("/sc-academic-projection-detail")
	{
		// Auth middleware menggunakan userService (bukan detailService)
		detailGroup.Use(
			middleware.AuthMiddleware(userService), // Diperbaiki: Gunakan userService seperti pada user
			middleware.RoleMiddleware("admin", "user"),
		)

		detailGroup.POST("", detailController.CreateScAcademicProjectionDetail)
		detailGroup.POST("/multiple", detailController.CreateMultipleScAcademicProjectionDetail) // Untuk foreach
		detailGroup.GET("", detailController.GetAllScAcademicProjectionDetail)
		detailGroup.GET("/:id", detailController.GetScAcademicProjectionDetailByID)
		detailGroup.PUT("/:id", detailController.UpdateScAcademicProjectionDetail)
		detailGroup.DELETE("/:id", detailController.DeleteScAcademicProjectionDetail)
	}
}
