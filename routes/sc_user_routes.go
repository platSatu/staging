// internal/routes/sc_user_routes.go

package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitSCUserRoutes(r *gin.Engine, db *gorm.DB) {
	userService := service.NewUserService(db)
	scUserService := service.NewSCUserService(db)
	scUserController := controller.NewSCUserController(scUserService)

	scUserGroup := r.Group("/sc-user")
	{
		// Auth middleware
		scUserGroup.Use(middleware.AuthMiddleware(userService))

		// Routes
		scUserGroup.POST("", scUserController.CreateSCUser)
		scUserGroup.GET("", scUserController.GetAllSCUsers)
		scUserGroup.GET("/my", scUserController.GetSCUsersByParentID)

		scUserGroup.GET("/:id", scUserController.GetSCUserByID)
		scUserGroup.PUT("/:id", scUserController.UpdateSCUser)
		scUserGroup.DELETE("/:id", scUserController.DeleteSCUser)
	}
}
