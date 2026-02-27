package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitScHabitsRoutes(r *gin.Engine, db *gorm.DB) {
	// Buat userService untuk middleware (konsisten dengan user_routes.go)
	userService := service.NewUserService(db)

	habitsService := service.NewScHabitsService(db)
	habitsController := controller.NewScHabitsController(habitsService)

	habitsGroup := r.Group("/sc-habits")
	{
		// Auth middleware menggunakan userService (bukan habitsService)
		habitsGroup.Use(
			middleware.AuthMiddleware(userService), // Diperbaiki: Gunakan userService seperti pada user
			middleware.RoleMiddleware("admin", "user"),
		)

		habitsGroup.POST("", habitsController.CreateScHabits)
		habitsGroup.GET("", habitsController.GetAllScHabits)
		habitsGroup.GET("/:id", habitsController.GetScHabitsByID)
		habitsGroup.PUT("/:id", habitsController.UpdateScHabits)
		habitsGroup.DELETE("/:id", habitsController.DeleteScHabits)
	}
}