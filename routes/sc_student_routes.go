package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitScStudentRoutes(r *gin.Engine, db *gorm.DB) {
	// Buat userService untuk middleware (konsisten dengan user_routes.go)
	userService := service.NewUserService(db)

	studentService := service.NewScStudentService(db)
	studentController := controller.NewScStudentController(studentService)

	studentGroup := r.Group("/sc-student")
	{
		// Auth middleware menggunakan userService (bukan studentService)
		studentGroup.Use(
			middleware.AuthMiddleware(userService), // Diperbaiki: Gunakan userService seperti pada user
			middleware.RoleMiddleware("admin", "user"),
		)

		studentGroup.POST("", studentController.CreateScStudent)
		studentGroup.GET("", studentController.GetAllScStudent)
		studentGroup.GET("/:id", studentController.GetScStudentByID)
		studentGroup.PUT("/:id", studentController.UpdateScStudent)
		studentGroup.DELETE("/:id", studentController.DeleteScStudent)
	}
}
