package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitScSubHabitsRoutes(r *gin.Engine, db *gorm.DB) {
	// Buat userService untuk middleware (konsisten dengan user_routes.go)
	userService := service.NewUserService(db)

	subHabitsService := service.NewScSubHabitsService(db)
	subHabitsController := controller.NewScSubHabitsController(subHabitsService)

	subHabitsGroup := r.Group("/sc-sub-habits")
	{
		// Auth middleware menggunakan userService (bukan subHabitsService)
		subHabitsGroup.Use(
			middleware.AuthMiddleware(userService), // Diperbaiki: Gunakan userService seperti pada user
			middleware.RoleMiddleware("admin", "user"),
		)

		subHabitsGroup.POST("", subHabitsController.CreateScSubHabits)
		subHabitsGroup.GET("", subHabitsController.GetAllScSubHabits)
		subHabitsGroup.GET("/:id", subHabitsController.GetScSubHabitsByID)
		subHabitsGroup.PUT("/:id", subHabitsController.UpdateScSubHabits)
		subHabitsGroup.DELETE("/:id", subHabitsController.DeleteScSubHabits)
	}
}