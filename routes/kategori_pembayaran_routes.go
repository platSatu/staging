package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitKategoriPembayaranRoutes(r *gin.Engine, db *gorm.DB) {
	// Buat userService untuk middleware (konsisten dengan user_routes.go)
	userService := service.NewUserService(db)

	kategoriService := service.NewKategoriPembayaranService(db)
	kategoriController := controller.NewKategoriPembayaranController(kategoriService)

	kategoriGroup := r.Group("/kategori-pembayaran")
	{
		// Auth middleware menggunakan userService (bukan kategoriService)
		kategoriGroup.Use(
			middleware.AuthMiddleware(userService), // Diperbaiki: Gunakan userService seperti pada user
			middleware.RoleMiddleware("admin", "user"),
		)

		kategoriGroup.POST("", kategoriController.CreateKategoriPembayaran)
		kategoriGroup.GET("", kategoriController.GetAllKategoriPembayaran)
		kategoriGroup.GET("/:id", kategoriController.GetKategoriPembayaranByID)
		kategoriGroup.PUT("/:id", kategoriController.UpdateKategoriPembayaran)
		kategoriGroup.DELETE("/:id", kategoriController.DeleteKategoriPembayaran)
	}
}
