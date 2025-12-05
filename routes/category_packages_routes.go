package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitCategoryPackagesRoutes(r *gin.Engine, db *gorm.DB) {
	categoryService := service.NewCategoryPackagesService(db)
	categoryController := controller.NewCategoryPackagesController(categoryService)

	categoryGroup := r.Group("/category_packages")
	{
		categoryGroup.Use(
			middleware.AuthMiddleware(service.NewUserService(db)), // Menggunakan UserService untuk auth
			middleware.RoleMiddleware("admin", "user"),
		)

		categoryGroup.POST("", categoryController.CreateCategoryPackages)
		categoryGroup.GET("", categoryController.GetAllCategoryPackages)
		categoryGroup.GET("/:id", categoryController.GetCategoryPackagesByID)
		categoryGroup.GET("/user", categoryController.GetCategoryPackagesByUser) // Endpoint khusus untuk user login
		categoryGroup.PUT("/:id", categoryController.UpdateCategoryPackages)
		categoryGroup.DELETE("/:id", categoryController.DeleteCategoryPackages)
	}
}
