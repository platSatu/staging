package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitScTemplateHabitsRoutes(r *gin.Engine, db *gorm.DB) {
	// Buat userService untuk middleware (konsisten dengan user_routes.go)
	userService := service.NewUserService(db)

	templateService := service.NewScTemplateHabitsService(db)
	templateController := controller.NewScTemplateHabitsController(templateService)

	templateGroup := r.Group("/sc-template-habits")
	{
		// Auth middleware menggunakan userService (bukan templateService)
		templateGroup.Use(
			middleware.AuthMiddleware(userService), // Diperbaiki: Gunakan userService seperti pada user
			middleware.RoleMiddleware("admin", "user"),
		)

		templateGroup.POST("", templateController.CreateScTemplateHabits)
		templateGroup.GET("", templateController.GetAllScTemplateHabits)
		templateGroup.GET("/:id", templateController.GetScTemplateHabitsByID)
		templateGroup.PUT("/:id", templateController.UpdateScTemplateHabits)
		templateGroup.DELETE("/:id", templateController.DeleteScTemplateHabits)
	}
}