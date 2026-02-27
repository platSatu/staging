package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitScGradeRoutes(r *gin.Engine, db *gorm.DB) {
	// Buat userService untuk middleware (konsisten dengan user_routes.go)
	userService := service.NewUserService(db)

	gradeService := service.NewScGradeService(db)
	gradeController := controller.NewScGradeController(gradeService)

	gradeGroup := r.Group("/sc-grade")
	{
		// Auth middleware menggunakan userService (bukan gradeService)
		gradeGroup.Use(
			middleware.AuthMiddleware(userService), // Diperbaiki: Gunakan userService seperti pada user
			middleware.RoleMiddleware("admin", "user"),
		)

		gradeGroup.POST("", gradeController.CreateScGrade)
		gradeGroup.GET("", gradeController.GetAllScGrade)
		gradeGroup.GET("/:id", gradeController.GetScGradeByID)
		gradeGroup.PUT("/:id", gradeController.UpdateScGrade)
		gradeGroup.DELETE("/:id", gradeController.DeleteScGrade)
	}
}