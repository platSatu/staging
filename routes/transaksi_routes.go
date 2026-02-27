package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitTransaksiRoutes(r *gin.Engine, db *gorm.DB) {
	// Buat userService untuk middleware
	userService := service.NewUserService(db)

	transaksiService := service.NewTransaksiService(db)
	transaksiController := controller.NewTransaksiController(transaksiService)

	transaksiGroup := r.Group("/transaksi")
	{
		// Auth middleware menggunakan userService
		transaksiGroup.Use(
			middleware.AuthMiddleware(userService),
			middleware.RoleMiddleware("admin", "user"),
		)
		transaksiGroup.GET("/user", transaksiController.GetTransaksiByUser)
		transaksiGroup.POST("", transaksiController.CreateTransaksi)
		transaksiGroup.GET("", transaksiController.GetAllTransaksi)
		transaksiGroup.GET("/:id", transaksiController.GetTransaksiByID)
		transaksiGroup.PUT("/:id", transaksiController.UpdateTransaksi)
		transaksiGroup.DELETE("/:id", transaksiController.DeleteTransaksi)
	}
}
