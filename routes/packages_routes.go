package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitPackagesRoutes(r *gin.Engine, db *gorm.DB) {
	packagesService := service.NewPackagesService(db)
	packagesController := controller.NewPackagesController(packagesService)

	packagesGroup := r.Group("/packages")
	{
		packagesGroup.Use(
			middleware.AuthMiddleware(service.NewUserService(db)), // Menggunakan UserService untuk auth
			middleware.RoleMiddleware("admin", "user"),
		)

		packagesGroup.POST("", packagesController.CreatePackages)
		packagesGroup.GET("", packagesController.GetAllPackages)
		packagesGroup.GET("/:id", packagesController.GetPackagesByID)
		packagesGroup.GET("/user", packagesController.GetPackagesByUser) // Endpoint khusus untuk user login
		packagesGroup.PUT("/:id", packagesController.UpdatePackages)
		packagesGroup.DELETE("/:id", packagesController.DeletePackages)
	}
}
