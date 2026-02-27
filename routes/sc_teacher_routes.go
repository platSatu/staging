package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitScTeacherRoutes(r *gin.Engine, db *gorm.DB) {
	// Buat userService untuk middleware (konsisten dengan user_routes.go)
	userService := service.NewUserService(db)

	teacherService := service.NewScTeacherService(db)
	teacherController := controller.NewScTeacherController(teacherService)

	teacherGroup := r.Group("/sc-teacher")
	{
		// Auth middleware menggunakan userService (bukan teacherService)
		teacherGroup.Use(
			middleware.AuthMiddleware(userService), // Diperbaiki: Gunakan userService seperti pada user
			middleware.RoleMiddleware("admin", "user"),
		)

		teacherGroup.POST("", teacherController.CreateScTeacher)
		teacherGroup.GET("", teacherController.GetAllScTeacher)
		teacherGroup.GET("/:id", teacherController.GetScTeacherByID)
		teacherGroup.PUT("/:id", teacherController.UpdateScTeacher)
		teacherGroup.DELETE("/:id", teacherController.DeleteScTeacher)
	}
}
