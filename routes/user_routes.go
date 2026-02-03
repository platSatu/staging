package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitUserRoutes(r *gin.Engine, db *gorm.DB) {
	userService := service.NewUserService(db)
	userController := controller.NewUserController(userService)

	userGroup := r.Group("/users")
	{
		// Auth middleware untuk semua endpoint user
		userGroup.Use(
			middleware.AuthMiddleware(userService),
			middleware.RoleMiddleware("admin", "user", "admin_events"),
		)

		// 🔥 PENTING: jangan pakai slash "/" di akhir!
		userGroup.POST("", userController.CreateUser)
		userGroup.GET("", userController.GetAllUsers)
		userGroup.GET("/:id", userController.GetUserByID)
		userGroup.PUT("/:id", userController.UpdateUser)
		userGroup.DELETE("/:id", userController.DeleteUser)

		// endpoint profile
		userGroup.GET("/profile", userController.GetProfile)
		userGroup.GET("/children", userController.GetChildren)
	}
}
