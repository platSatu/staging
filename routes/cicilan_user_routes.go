package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitCicilanUserRoutes(r *gin.Engine, db *gorm.DB) {
	// Buat userService untuk middleware
	userService := service.NewUserService(db)

	cicilanService := service.NewCicilanUserService(db)
	cicilanController := controller.NewCicilanUserController(cicilanService)

	cicilanGroup := r.Group("/cicilan-user")
	{
		// Auth middleware menggunakan userService
		cicilanGroup.Use(
			middleware.AuthMiddleware(userService),
			middleware.RoleMiddleware("admin", "user"),
		)

		cicilanGroup.POST("", cicilanController.CreateCicilanUser)
		cicilanGroup.GET("", cicilanController.GetAllCicilanUser)
		cicilanGroup.GET("/:id", cicilanController.GetCicilanUserByID)
		cicilanGroup.PUT("/:id", cicilanController.UpdateCicilanUser)
		cicilanGroup.DELETE("/:id", cicilanController.DeleteCicilanUser)

		cicilanGroup.GET("/parent", cicilanController.GetParentList)
		cicilanGroup.GET("/parent/:parent_id", cicilanController.GetCicilanByParentID)
	}
}
