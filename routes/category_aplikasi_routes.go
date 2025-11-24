package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitCategoryAplikasiRoutes(r *gin.Engine, db *gorm.DB) {
	categoryService := service.NewCategoryAplikasiService(db)
	categoryController := controller.NewCategoryAplikasiController(categoryService)

	categoryGroup := r.Group("/category_aplikasi")
	{
		categoryGroup.Use(
			middleware.AuthMiddleware(service.NewUserService(db)), // Menggunakan UserService untuk auth
			middleware.RoleMiddleware("admin", "user"),
		)

		categoryGroup.POST("", categoryController.CreateCategoryAplikasi)
		categoryGroup.GET("", categoryController.GetAllCategoryAplikasi)
		categoryGroup.GET("/:id", categoryController.GetCategoryAplikasiByID)
		categoryGroup.GET("/user", categoryController.GetCategoryAplikasiByUser) // Endpoint khusus untuk user login
		categoryGroup.PUT("/:id", categoryController.UpdateCategoryAplikasi)
		categoryGroup.DELETE("/:id", categoryController.DeleteCategoryAplikasi)
	}
}
