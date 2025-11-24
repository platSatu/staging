package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// InitAuthRoutes fokus ke autentikasi (register, login, logout)
func InitAuthRoutes(r *gin.Engine, db *gorm.DB) {
	// Buat service
	userService := service.NewUserService(db)
	authService := service.NewAuthService(db, userService) // sesuai signature

	// Buat controller
	authController := controller.NewAuthController(authService)

	authGroup := r.Group("/auth") // prefix /auth
	{
		authGroup.POST("/register", authController.Register)
		authGroup.POST("/login", authController.Login)
		authGroup.POST("/logout", authController.Logout)
		authGroup.POST("/refresh", authController.Refresh)
	}
}
