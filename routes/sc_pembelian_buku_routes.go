package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitScPembelianBukuRoutes(r *gin.Engine, db *gorm.DB) {
	// Buat userService untuk middleware
	userService := service.NewUserService(db)

	pembelianBukuService := service.NewScPembelianBukuService(db)
	pembelianBukuController := controller.NewScPembelianBukuController(pembelianBukuService)

	pembelianBukuGroup := r.Group("/sc-pembelian-buku")
	{
		// Auth middleware menggunakan userService
		pembelianBukuGroup.Use(
			middleware.AuthMiddleware(userService),
			middleware.RoleMiddleware("admin", "user"),
		)

		pembelianBukuGroup.POST("", pembelianBukuController.CreateScPembelianBuku)
		pembelianBukuGroup.GET("", pembelianBukuController.GetAllScPembelianBuku)
		pembelianBukuGroup.GET("/:id", pembelianBukuController.GetScPembelianBukuByID)
		pembelianBukuGroup.PUT("/:id", pembelianBukuController.UpdateScPembelianBuku)
		pembelianBukuGroup.DELETE("/:id", pembelianBukuController.DeleteScPembelianBuku)
	}
}
