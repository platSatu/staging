package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitLaporanRoutes(r *gin.Engine, db *gorm.DB) {

	// userService untuk middleware auth
	userService := service.NewUserService(db)

	// layanan laporan
	laporanService := service.NewLaporanService(db)
	laporanController := controller.NewLaporanController(laporanService)

	laporanGroup := r.Group("/laporan")
	{
		laporanGroup.Use(
			middleware.AuthMiddleware(userService),
			middleware.RoleMiddleware("admin", "user"),
		)

		// Laporan per kategori
		laporanGroup.GET("/perkategori", laporanController.GetLaporanPerKategori)
		laporanGroup.GET("/form/:id", laporanController.GetDetailPerForm)

	}
}
