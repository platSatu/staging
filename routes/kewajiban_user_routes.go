package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitKewajibanUserRoutes(r *gin.Engine, db *gorm.DB) {
	// Buat userService untuk middleware
	userService := service.NewUserService(db)

	kewajibanService := service.NewKewajibanUserService(db)
	kewajibanController := controller.NewKewajibanUserController(kewajibanService)

	kewajibanGroup := r.Group("/kewajiban-user")
	{
		// Auth middleware menggunakan userService
		kewajibanGroup.Use(
			middleware.AuthMiddleware(userService),
			middleware.RoleMiddleware("admin", "user"),
		)

		kewajibanGroup.POST("", kewajibanController.CreateKewajibanUser)
		kewajibanGroup.GET("", kewajibanController.GetAllKewajibanUser)
		kewajibanGroup.GET("/:id", kewajibanController.GetKewajibanUserByID)
		kewajibanGroup.PUT("/:id", kewajibanController.UpdateKewajibanUser)
		kewajibanGroup.DELETE("/:id", kewajibanController.DeleteKewajibanUser)
	}
}
